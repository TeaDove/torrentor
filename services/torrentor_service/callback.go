package torrentor_service

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
	"torrentor/schemas"
	"torrentor/services/ffmpeg_service"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) onTorrentComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
) error {
	<-torrentEnt.Obj.Complete().On()

	torrentEnt.Completed = true

	zerolog.Ctx(ctx).Info().Dict("torrent", torrentEnt.ZerologDict()).Msg("torrent.ready")

	return nil
}

func (r *Service) onFileCompleteCallback(
	ctx context.Context,
	fileEnt *schemas.FileEntityPop,
) error {
	var err error

	if fileEnt.Mimetype == schemas.MatroskaMimeType {
		err = r.unpackMatroska(ctx, fileEnt)
		if err != nil {
			return errors.Wrap(err, "failed to unpack matroska file")
		}
	}

	fileEnt.Completed = true

	zerolog.Ctx(ctx).
		Trace().
		Interface("file", fileEnt.Name).
		Msg("file.ready")

	return nil
}

func (r *Service) onFileComplete(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntityPop,
	completedCheckPeriod time.Duration,
) error {
	// TODO check if already completed
	incompleteFiles := map[string]*torrent.File{}
	for _, file := range torrentEnt.Obj.Files() {
		incompleteFiles[file.Path()] = file
	}

	completed := make([]string, 0, len(incompleteFiles))

	for {
		for _, file := range incompleteFiles {
			if file.Length() == file.BytesCompleted() {
				completed = append(completed, file.Path())
			}
		}

		for _, fileName := range completed {
			err := r.onFileCompleteCallback(ctx, &schemas.FileEntityPop{
				FileEntity: torrentEnt.Files[schemas.TrimFirstDir(fileName)],
				Obj:        incompleteFiles[fileName],
				Torrent:    torrentEnt,
			})
			if err != nil {
				return errors.Wrap(err, "failed to unpack matroska")
			}

			delete(incompleteFiles, fileName)
		}

		completed = make([]string, 0, len(incompleteFiles))

		if len(incompleteFiles) == 0 {
			break
		}

		time.Sleep(completedCheckPeriod)
	}

	err := r.onTorrentComplete(ctx, torrentEnt)
	if err != nil {
		return errors.Wrap(err, "failed to mark complete")
	}

	return nil
}

func makeFilenameWithTags(base string, ext string, tags ...string) string {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		base += "-" + tag
	}

	return fmt.Sprintf("%s%s", base, ext)
}

func (r *Service) saveFile(
	ctx context.Context,
	torrentEnt *schemas.TorrentEntity,
	newFilePath string,
) error {
	fileStats, err := os.Stat(newFilePath)
	if err != nil {
		return errors.Wrap(err, "error opening file")
	}

	newFileEnt := makeFileEnt(
		schemas.TrimFirstDir(schemas.TrimFirstDir(schemas.TrimFirstDir(newFilePath))),
		uint64(fileStats.Size()),
		true,
	)

	torrentEnt.AppendFile(newFileEnt)
	torrentEnt.Files[newFileEnt.Path] = newFileEnt

	return nil
}

func (r *Service) unpackMatroska(
	ctx context.Context,
	fileEnt *schemas.FileEntityPop,
) error {
	filePath := fileEnt.Location(r.torrentDataDir)
	newFilesDir := filepath.Join(filepath.Dir(filePath), fileEnt.NameWithoutExt())
	err := os.MkdirAll(newFilesDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error creating file directory")
	}

	metadata, err := r.ffmpegService.ExportMetadata(ctx, filePath)
	if err != nil {
		return errors.Wrap(err, "error exporting metadata")
	}

	audioIdx := 0
	subIdx := 0

	var newFilename string

	for _, stream := range metadata.Streams {
		switch stream.CodecType {
		case ffmpeg_service.CodecTypeAudio:
			newFilename = path.Join(newFilesDir, makeFilenameWithTags(
				fileEnt.NameWithoutExt(),
				".mp4",
				stream.Tags.Title,
				stream.Tags.Language,
			))

			err = r.ffmpegService.MKVExportMP4(ctx, filePath, audioIdx, newFilename)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			newFolder := path.Join(newFilesDir, makeFilenameWithTags(
				fileEnt.NameWithoutExt(),
				"/hls/",
				stream.Tags.Title,
				stream.Tags.Language,
			))

			err = r.ffmpegService.MKVExportHLS(ctx, newFilename, audioIdx, newFolder)
			if err != nil {
				return errors.Wrap(err, "error converting mp4 to hls")
			}

			audioIdx++
		case ffmpeg_service.CodecTypeSubtitle:
			newFilename = path.Join(newFilesDir, makeFilenameWithTags(
				fileEnt.NameWithoutExt(),
				".vtt",
				stream.Tags.Title,
				stream.Tags.Language,
			))

			err = r.ffmpegService.MKVExportSubtitles(ctx, filePath, subIdx, newFilename)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			subIdx++
		default:
			continue
		}

		err = r.saveFile(ctx, fileEnt.Torrent.TorrentEntity, newFilename)
		if err != nil {
			return errors.Wrap(err, "error adding file to DB")
		}
	}

	return nil
}

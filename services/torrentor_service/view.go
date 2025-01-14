package torrentor_service

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"torrentor/schemas"
	"torrentor/services/ffmpeg_service"
	"torrentor/settings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (schemas.TorrentEntity, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

func (r *Service) GetFileWithContent(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (schemas.FileWithContent, error) {
	torrent, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "error getting torrent")
	}

	file, err := os.Open(torrent.FileLocation(settings.Settings.Torrent.DataDir, filePath))
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "error opening file")
	}

	fileMeta, ok := torrent.Files[filePath]
	if !ok {
		return schemas.FileWithContent{}, errors.New("file not found")
	}

	return schemas.FileWithContent{FileEntity: fileMeta, OSFile: file}, nil
}

func (r *Service) GetFile(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (schemas.FileEntity, error) {
	torrentEnt, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return schemas.FileEntity{}, errors.Wrap(err, "error getting torrent")
	}

	file, ok := torrentEnt.Files[filePath]
	if !ok {
		return schemas.FileEntity{}, errors.New("file not found")
	}

	// if file.Mimetype == schemas.MatroskaMimeType {
	//	err = r.unpackMatroska(ctx, torrentEnt.FileLocation(settings.Settings.Torrent.DataDir, filePath))
	//	if err != nil {
	//		return schemas.FileEntity{}, errors.Wrap(err, "error unpacking file")
	//	}
	//}

	return file, nil
}

func makeFilenameWithTags(base string, ext string, tags ...string) string {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		base += "-" + tag
	}

	return fmt.Sprintf("%s.%s", base, ext)
}

func (r *Service) addFileToDB(
	ctx context.Context,
	fileEntOriginal *schemas.FileEntityPop,
	newFilePath string,
) error {
	fileStats, err := os.Stat(newFilePath)
	if err != nil {
		return errors.Wrap(err, "error opening file")
	}

	oldFileDir := filepath.Dir(fileEntOriginal.Path)
	newFileBase := filepath.Base(newFilePath)

	newFileEnt := makeFileEnt(filepath.Join("a", oldFileDir, newFileBase), uint64(fileStats.Size()), true)
	fileEntOriginal.Torrent.AppendFile(newFileEnt)

	_, err = r.torrentRepository.TorrentSave(ctx, &fileEntOriginal.Torrent.TorrentEntity)
	if err != nil {
		return errors.Wrap(err, "error upserting torrent")
	}

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
				"mp4",
				stream.Tags.Title,
				stream.Tags.Language,
			))

			err = r.ffmpegService.MKVExportMP4(ctx, filePath, audioIdx, newFilename)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			audioIdx++
		case ffmpeg_service.CodecTypeSubtitle:
			newFilename = path.Join(newFilesDir, makeFilenameWithTags(
				fileEnt.NameWithoutExt(),
				"vtt",
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

		err = r.addFileToDB(ctx, fileEnt, newFilename)
		if err != nil {
			return errors.Wrap(err, "error adding file to DB")
		}
	}

	return nil
}

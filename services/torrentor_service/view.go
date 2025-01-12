package torrentor_service

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"torrentor/repositories/torrent_repository"
	"torrentor/services/ffmpeg_service"
	"torrentor/settings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (torrent_repository.Torrent, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

func (r *Service) GetFileWithContent(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (torrent_repository.FileWithContent, error) {
	torrent, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return torrent_repository.FileWithContent{}, errors.Wrap(err, "error getting torrent")
	}

	file, err := os.Open(torrent.FileLocation(settings.Settings.Torrent.DataDir, filePath))
	if err != nil {
		return torrent_repository.FileWithContent{}, errors.Wrap(err, "error opening file")
	}

	fileMeta, ok := torrent.Files[filePath]
	if !ok {
		return torrent_repository.FileWithContent{}, errors.New("file not found")
	}

	return torrent_repository.FileWithContent{File: fileMeta, OSFile: file}, nil
}

func (r *Service) GetFile(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (torrent_repository.File, error) {
	torrent, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return torrent_repository.File{}, errors.Wrap(err, "error getting torrent")
	}

	file, ok := torrent.Files[filePath]
	if !ok {
		return torrent_repository.File{}, errors.New("file not found")
	}

	if file.Mimetype == torrent_repository.MatroskaMimeType {
		err = r.unpackMatroska(ctx, torrent.FileLocation(settings.Settings.Torrent.DataDir, filePath))
		if err != nil {
			return torrent_repository.File{}, errors.Wrap(err, "error unpacking file")
		}
	}

	return file, nil
}

func (r *Service) unpackMatroska(ctx context.Context, filePath string) error {
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	metadata, err := r.ffmpegService.ExportMetadata(ctx, filePath)
	if err != nil {
		return errors.Wrap(err, "error exporting metadata")
	}

	audioIdx := 0
	subIdx := 0

	for _, stream := range metadata.Streams {
		switch stream.CodecType {
		case ffmpeg_service.CodecTypeAudio:
			err = r.ffmpegService.MKVExportMP4(
				ctx,
				filePath,
				audioIdx,
				path.Join(path.Dir(filePath), fmt.Sprintf("%s-%s-%s.mp4",
					fileNameWithoutExt,
					stream.Tags.Title,
					stream.Tags.Language,
				)),
			)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			audioIdx++
		case ffmpeg_service.CodecTypeSubtitle:
			err = r.ffmpegService.MKVExportSubtitles(
				ctx,
				filePath,
				subIdx,
				path.Join(path.Dir(filePath), fmt.Sprintf("%s-%s-%s.vtt",
					fileNameWithoutExt,
					stream.Tags.Title,
					stream.Tags.Language,
				)),
			)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			subIdx++
		default:
			continue
		}
	}

	return nil
}

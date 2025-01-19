package torrentor_service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"torrentor/schemas"
	"torrentor/services/ffmpeg_service"
)

func makeFilenameWithTags(base string, suffix string, tags ...string) string {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		base += "-" + tag
	}

	return fmt.Sprintf("%s%s", base, suffix)
}

//func (r *Service) saveFile(
//	ctx context.Context,
//	torrentEnt *schemas.TorrentEntity,
//	newFilePath string,
//) error {
//	fileStats, err := os.Stat(newFilePath)
//	if err != nil {
//		return errors.Wrap(err, "error opening file")
//	}
//
//	newFileEnt := makeFileEnt(
//		schemas.TrimFirstDir(schemas.TrimFirstDir(schemas.TrimFirstDir(newFilePath))),
//		uint64(fileStats.Size()),
//		nil,
//		true,
//	)
//
//	torrentEnt.AppendFile(newFileEnt)
//	torrentEnt.FilePathMap[newFileEnt.Path] = newFileEnt
//
//	return nil
//}

//func (r *Service) matroskaToHLS(
//	ctx context.Context,
//	file string,
//) (string, error) {
//}

//func (r *Service) unpackMatroskaAudio(
//	ctx context.Context,
//	unpackFilesDir string,
//	fileEnt *schemas.FileEntityPop,
//	stream ffmpeg_service.Stream,
//) error {
//	newFilename := path.Join(unpackFilesDir, makeFilenameWithTags(
//		fileEnt.Name,
//		".mp4",
//		stream.Tags.Title,
//		stream.Tags.Language,
//	))
//
//	err = r.ffmpegService.MKVExportMP4(ctx, filePath, audioIdx, newFilename)
//	if err != nil {
//		return errors.Wrap(err, "error converting audio stream")
//	}
//
//	newFolder := path.Join(unpackFilesDir, makeFilenameWithTags(
//		"hls",
//		"/",
//		stream.Tags.Title,
//		stream.Tags.Language,
//	), "output.m3u8")
//
//	err = r.ffmpegService.MKVExportHLS(ctx, newFilename, audioIdx, newFolder)
//	if err != nil {
//		return errors.Wrap(err, "error converting mp4 to hls")
//	}
//}

func (r *Service) unpackMatroska(
	ctx context.Context,
	fileEnt *schemas.FileEntity,
) error {
	filePath := fileEnt.Location()

	unpackFilesDir := filepath.Dir(fileEnt.Location())
	err := os.MkdirAll(unpackFilesDir, os.ModePerm)
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
			newFilename = path.Join(unpackFilesDir, makeFilenameWithTags(
				fileEnt.Name,
				".mp4",
				stream.Tags.Title,
				stream.Tags.Language,
			))

			err = r.ffmpegService.MKVExportMP4(ctx, filePath, audioIdx, newFilename)
			if err != nil {
				return errors.Wrap(err, "error converting audio stream")
			}

			newFolder := path.Join(unpackFilesDir, makeFilenameWithTags(
				"hls",
				"/",
				stream.Tags.Title,
				stream.Tags.Language,
			), "output.m3u8")

			err = r.ffmpegService.MKVExportHLS(ctx, newFilename, audioIdx, newFolder)
			if err != nil {
				return errors.Wrap(err, "error converting mp4 to hls")
			}

			audioIdx++
		case ffmpeg_service.CodecTypeSubtitle:
			newFilename = path.Join(unpackFilesDir, makeFilenameWithTags(
				fileEnt.Name,
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

		//err = r.saveFile(ctx, fileEnt.Torrent.TorrentEntity, newFilename)
		//if err != nil {
		//	return errors.Wrap(err, "error adding file to DB")
		//}
	}

	return nil
}

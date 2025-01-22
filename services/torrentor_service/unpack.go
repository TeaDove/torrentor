package torrentor_service

import (
	"context"
	"os"
	"torrentor/schemas"
	"torrentor/services/ffmpeg_service"

	"github.com/pkg/errors"
)

func (r *Service) unpackMatroskaAudio(
	ctx context.Context,
	mkvFileEnt *schemas.FileEntity,
	stream *ffmpeg_service.Stream,
	audioIdx int,
) error {
	mp4File := mkvFileEnt.LocationInUnpackAsStream(stream, ".mp4")

	hlsFolder := mkvFileEnt.LocationInUnpackAsStream(stream, ".m3u8/output.m3u8")
	if _, err := os.Stat(hlsFolder); err == nil {
		return nil
	}

	err := r.ffmpegService.MKVExportMP4(ctx, mkvFileEnt.Location(), audioIdx, mp4File)
	if err != nil {
		return errors.Wrap(err, "error converting audio stream")
	}

	err = r.ffmpegService.MKVExportHLS(ctx, mp4File, 0, hlsFolder)
	if err != nil {
		return errors.Wrap(err, "error converting mp4 to hls")
	}

	err = os.Remove(mp4File)
	if err != nil {
		return errors.Wrap(err, "error removing file")
	}

	return nil
}

func (r *Service) unpackMatroska(
	ctx context.Context,
	mkvFileEnt *schemas.FileEntity,
) error {
	filePath := mkvFileEnt.Location()

	err := os.MkdirAll(mkvFileEnt.LocationInUnpack(), os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "error creating file directory")
	}

	metadata, err := r.ffmpegService.ExportMetadata(ctx, filePath)
	if err != nil {
		return errors.Wrap(err, "error exporting metadata")
	}

	audioIdx := 0
	subIdx := 0

	for _, stream := range metadata.Streams {
		switch stream.CodecType {
		case ffmpeg_service.CodecTypeAudio:
			err = r.unpackMatroskaAudio(ctx, mkvFileEnt, &stream, audioIdx)
			if err != nil {
				return errors.Wrap(err, "failed to unpack audio stream")
			}

			audioIdx++
		case ffmpeg_service.CodecTypeSubtitle:
			err = r.ffmpegService.MKVExportSubtitles(
				ctx,
				filePath,
				subIdx,
				mkvFileEnt.LocationInUnpackAsStream(&stream, ".vtt"),
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

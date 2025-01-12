package ffmpeg_service

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
)

func (r *Service) MKVUnpack(ctx context.Context, filePath string, writer io.Writer) error {
	var ffmpegErrBuf bytes.Buffer

	metadataObj, err := r.exportMetadata(ctx, filePath)
	if err != nil {
		return errors.Wrap(err, "failed to export metadata")
	}

	for _, streamObj := range metadataObj.Streams {
		var (
			fileName string
			outKW    = ffmpeg.KwArgs{}
			inKW     = ffmpeg.KwArgs{}
		)
		switch streamObj.CodecType {
		case codecTypeSubtitle:
			fileName = "output.vtt"
			outKW["c:s"] = "webvtt"
		case codecTypeAudio:
			fileName = "output.mp3"
			outKW["map"] = "0:a"
			outKW["codec"] = "copy"
		case codecTypeVideo:
			fileName = "output.mp4"
			outKW["map"] = "0:v"
			outKW["codec"] = "copy"
		default:
			zerolog.Ctx(ctx).
				Warn().
				Str("file", filePath).
				Interface("metadata", metadataObj).
				Msg("unknown.codec.type")
			continue
		}

		err = ffmpeg.
			Input(filePath, inKW).
			WithErrorOutput(&ffmpegErrBuf).
			Output(fileName, outKW).
			Run()

		if err != nil {
			if ffmpegErrBuf.Len() > 0 {
				return errors.Errorf("failed to run ffmpeg: %s", ffmpegErrBuf.String())
			}
			return errors.Wrap(err, "failed to run ffmpeg")
		}
	}

	return nil
}

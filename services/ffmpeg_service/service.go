package ffmpeg_service

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Service struct{}

func NewService(_ context.Context) (*Service, error) {
	return &Service{}, nil
}

func runWithErr(stream *ffmpeg.Stream) error {
	var ffmpegErrBuf bytes.Buffer

	err := stream.WithErrorOutput(&ffmpegErrBuf).OverWriteOutput().Run()
	if err != nil {
		if ffmpegErrBuf.Len() > 0 {
			return errors.Errorf("failed to run ffmpeg: %s", ffmpegErrBuf.String())
		}

		return errors.Wrap(err, "failed to run ffmpeg")
	}

	return nil
}

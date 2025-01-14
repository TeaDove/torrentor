package ffmpeg_service

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"strings"
	"time"
)

type Service struct{}

func NewService(_ context.Context) (*Service, error) {
	ffmpeg.LogCompiledCommand = false
	return &Service{}, nil
}

func runWithErr(ctx context.Context, stream *ffmpeg.Stream) error {
	var ffmpegErrBuf bytes.Buffer

	t0 := time.Now()
	compiledCommand := strings.Join(stream.Compile().Args, " ")

	err := stream.WithErrorOutput(&ffmpegErrBuf).OverWriteOutput().Run()
	zerolog.Ctx(ctx).
		Debug().
		Str("elapsed", time.Since(t0).String()).
		Str("command", compiledCommand).
		Msg("ffmpeg.called")

	if err != nil {
		if ffmpegErrBuf.Len() > 0 {
			return errors.Errorf("failed to run ffmpeg: %s", ffmpegErrBuf.String())
		}

		return errors.Wrap(err, "failed to run ffmpeg")
	}

	return nil
}

package ffmpeg_service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (r *Service) MKVGetSubtitles(ctx context.Context, idx int, distPath string) error {
	err := runWithErr(ffmpeg.
		Input(distPath).
		Output(distPath, ffmpeg.KwArgs{
			"c:s": "webvtt",
			"map": fmt.Sprintf("0:s:%d", idx),
		}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	return nil
}

func (r *Service) MKVToMP4(ctx context.Context, filePath string, audioIdx int, distPath string) error {
	err := runWithErr(ffmpeg.
		Input(filePath).
		Output(distPath,
			ffmpeg.KwArgs{"codec": "copy"},
			ffmpeg.KwArgs{"map": []string{"0:v:0", fmt.Sprintf("0:a:%d", audioIdx)}},
		))
	if err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	return nil
}

package ffmpeg_service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
)

func (r *Service) MKVExportSubtitles(ctx context.Context, filePath string, subIdx int, distPath string) error {
	if _, err := os.Stat(distPath); err == nil {
		return nil
	}

	err := runWithErr(ctx, ffmpeg.
		Input(filePath).
		Output(distPath, ffmpeg.KwArgs{
			"c:s": "webvtt",
			"map": fmt.Sprintf("0:s:%d", subIdx),
		}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	zerolog.Ctx(ctx).
		Info().
		Str("distPath", distPath).
		Msg("subtitles.exported")

	return nil
}

func (r *Service) MKVExportMP4(ctx context.Context, filePath string, audioIdx int, distPath string) error {
	if _, err := os.Stat(distPath); err == nil {
		return nil
	}

	err := runWithErr(ctx, ffmpeg.
		Input(filePath).
		Output(distPath,
			ffmpeg.KwArgs{"codec": "copy"},
			ffmpeg.KwArgs{"map": []string{"0:v:0", fmt.Sprintf("0:a:%d", audioIdx)}},
		))
	if err != nil {
		return errors.Wrap(err, "failed to run ffmpeg")
	}

	zerolog.Ctx(ctx).
		Info().
		Str("distPath", distPath).
		Msg("mp4.exported")

	return nil
}

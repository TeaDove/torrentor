package ffmpeg_service

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func TestUnit_FfmpegService_MKVUnpack_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	r, err := NewService(ctx)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = r.MKVUnpack(
		ctx,
		"/Users/pibragimov/projects/torrentor/data/torrent/0d7f1fe0531741902f8d6637ee787c99bff48791/Shameless.S03.720p.BDRip.x264.ac3.rus.eng/Shameless.S03.E01.BDRip.720p.mkv",
		&buf,
	)
	require.NoError(t, err)

	logger_utils.LogAny(buf.Len())
}

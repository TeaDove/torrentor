package ffmpeg_service

import (
	"testing"

	"github.com/teadove/teasutils/utils/test_utils"

	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func TestUnit_FfmpegService_MKVUnpack_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	r, err := NewService(ctx)
	require.NoError(t, err)

	err = r.MKVExportMP4(
		ctx,
		"/Users/pibragimov/projects/torrentor/backend/data/"+
			"torrent/0d7f1fe0531741902f8d6637ee787c99bff48791/"+
			"Shameless.S03.720p.BDRip.x264.ac3.rus.eng/Shameless.S03.E01.BDRip.720p.mkv",
		1,
		".test/output.mp4",
	)
	require.NoError(t, err)
}

func TestUnit_FfmpegService_Metadata_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	r, err := NewService(ctx)
	require.NoError(t, err)

	metadata, err := r.ExportMetadata(
		ctx,
		"/Users/pibragimov/projects/torrentor/backend/data/torrent/"+
			"0d7f1fe0531741902f8d6637ee787c99bff48791/Shameless.S03.720p.BDRip.x264.ac3.rus.eng/"+
			"Shameless.S03.E01.BDRip.720p.mkv",
	)
	require.NoError(t, err)
	test_utils.Pprint(metadata)
}

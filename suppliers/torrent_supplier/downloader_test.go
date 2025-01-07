package torrent_supplier

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func TestIntegration_TorrentSupplier_DownloadMagnetTo_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	expectedName := "Ремарк-Собрание в 16 томах-2011-EPUB"
	dataDir := "./.test/torrent_supplier"
	os.RemoveAll(path.Join(dataDir, expectedName))

	supplier, err := NewSupplier(ctx, dataDir)
	require.NoError(t, err)
	defer supplier.Close(ctx)

	//nolint: lll // i cant shorten it(
	actualPath, err := supplier.AddMagnetAndGetInfoAndStartDownload(
		ctx,
		"magnet:?xt=urn:btih:1AE80FD51FC9591C3369EC1BFA0EDBD3E6CDF019&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet&dn=%D0%AD%D1%80%D0%B8%D1%85%20%D0%9C%D0%B0%D1%80%D0%B8%D1%8F%20%D0%A0%D0%B5%D0%BC%D0%B0%D1%80%D0%BA%20-%20%D0%A1%D0%BE%D0%B1%D1%80%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%81%D0%BE%D1%87%D0%B8%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9%20%D0%B2%2016%20%D1%82%D0%BE%D0%BC%D0%B0%D1%85%20%5B2011%2C%20EPUB%2C%20RUS%5D",
	)
	require.NoError(t, err)
	assert.Equal(t, expectedName, actualPath)
	require.NoError(t, os.RemoveAll(path.Join(dataDir, expectedName)))
}

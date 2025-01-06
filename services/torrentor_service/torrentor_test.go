package torrentor_service

import (
	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/tidwall/buntdb"
	"testing"
	"torrentor/repositories/torrent_repository"
	"torrentor/suppliers/torrent_supplier"
)

func TestIntegration_TorrrentorService_SetGet_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	mangetLink := "magnet:?xt=urn:btih:1AE80FD51FC9591C3369EC1BFA0EDBD3E6CDF019&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet&dn=%D0%AD%D1%80%D0%B8%D1%85%20%D0%9C%D0%B0%D1%80%D0%B8%D1%8F%20%D0%A0%D0%B5%D0%BC%D0%B0%D1%80%D0%BA%20-%20%D0%A1%D0%BE%D0%B1%D1%80%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%81%D0%BE%D1%87%D0%B8%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9%20%D0%B2%2016%20%D1%82%D0%BE%D0%BC%D0%B0%D1%85%20%5B2011%2C%20EPUB%2C%20RUS%5D"

	supplier, err := torrent_supplier.NewSupplier(ctx, "./.test/torrent/")
	require.NoError(t, err)

	db, err := buntdb.Open("./.test/data.db")
	require.NoError(t, err)

	repository, err := torrent_repository.NewRepository(ctx, db)
	require.NoError(t, err)

	service, err := NewService(ctx, supplier, repository)
	require.NoError(t, err)

	torrent, err := service.DownloadAndSaveFromMagnet(ctx, mangetLink)
	require.NoError(t, err)

	gettedTorrent, err := service.torrentRepository.TorrentGetById(ctx, torrent.Id)
	require.NoError(t, err)

	require.Equal(t, torrent, gettedTorrent)
	require.NotEmpty(t, torrent.Name)
}

package torrentor_service

import (
	"testing"
	"torrentor/backend/schemas"

	"github.com/anacrolix/torrent/metainfo"

	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
)

func TestIntegration_TorrrentorService_SetGet_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	service := getService(ctx, t)
	mangetLink := "magnet:?xt=urn:btih:1AE80FD51FC9591C3369EC1BFA0EDBD3E6CDF019&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet&dn=%D0%AD%D1%80%D0%B8%D1%85%20%D0%9C%D0%B0%D1%80%D0%B8%D1%8F%20%D0%A0%D0%B5%D0%BC%D0%B0%D1%80%D0%BA%20-%20%D0%A1%D0%BE%D0%B1%D1%80%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%81%D0%BE%D1%87%D0%B8%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9%20%D0%B2%2016%20%D1%82%D0%BE%D0%BC%D0%B0%D1%85%20%5B2011%2C%20EPUB%2C%20RUS%5D"

	torrent, _, err := service.DownloadAndSaveFromMagnet(ctx, mangetLink)
	require.NoError(t, err)

	gettedTorrent, err := service.GetTorrentByInfoHash(ctx, torrent.InfoHash)
	require.NoError(t, err)

	require.Equal(t, torrent, gettedTorrent)
	require.NotEmpty(t, torrent.Name)
}

func TestIntegration_TorrrentorService_ConverMKV_Ok(t *testing.T) {
	t.Parallel()

	ctx := logger_utils.NewLoggedCtx()
	service := getService(ctx, t)

	_, err := service.GetAllTorrents(ctx)
	require.NoError(t, err)

	err = service.unpackMatroska(ctx, &schemas.FileEntity{
		Name:     "input.mkv",
		Path:     "input.mkv",
		Mimetype: schemas.MatroskaMimeType,
		Torrent: &schemas.TorrentEntity{
			FilePathMap: make(map[string]*schemas.FileEntity),
			InfoHash:    metainfo.NewHashFromHex("0d7f1fe0531741902f8d6637ee787c99bff48791"),
		},
	})
	require.NoError(t, err)
}

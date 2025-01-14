package torrentor_service

import (
	"context"
	"testing"
	"time"
	"torrentor/repositories/torrent_repository"
	"torrentor/services/ffmpeg_service"
	"torrentor/settings"
	"torrentor/suppliers/torrent_supplier"

	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/tidwall/buntdb"
)

func getService(ctx context.Context, t *testing.T) *Service {
	t.Helper()

	scheduler := gocron.NewScheduler(time.UTC)

	db, err := buntdb.Open(":memory:")
	require.NoError(t, err)

	torrentDataDir := ".test/torrent"

	torrentRepository, err := torrent_repository.NewRepository(ctx, db, torrentDataDir)
	require.NoError(t, err)

	torrentSupplier, err := torrent_supplier.NewSupplier(ctx, torrentDataDir)
	require.NoError(t, err)

	ffmpegService, err := ffmpeg_service.NewService(ctx)
	require.NoError(t, err)

	torrentorService, err := NewService(
		ctx,
		torrentSupplier,
		torrentRepository,
		ffmpegService,
		scheduler,
		settings.Settings.Torrent.DataDir,
	)
	require.NoError(t, err)

	scheduler.StartAsync()

	return torrentorService
}

func TestIntegration_TorrentorService_Job_Ok(t *testing.T) {
	ctx := logger_utils.NewLoggedCtx()
	service := getService(ctx, t)

	mangetLink := "magnet:?xt=urn:btih:1AE80FD51FC9591C3369EC1BFA0EDBD3E6CDF019&tr=http%3A%2F%2Fbt.t-ru.org%2Fann%3Fmagnet&dn=%D0%AD%D1%80%D0%B8%D1%85%20%D0%9C%D0%B0%D1%80%D0%B8%D1%8F%20%D0%A0%D0%B5%D0%BC%D0%B0%D1%80%D0%BA%20-%20%D0%A1%D0%BE%D0%B1%D1%80%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%81%D0%BE%D1%87%D0%B8%D0%BD%D0%B5%D0%BD%D0%B8%D0%B9%20%D0%B2%2016%20%D1%82%D0%BE%D0%BC%D0%B0%D1%85%20%5B2011%2C%20EPUB%2C%20RUS%5D"

	_, _, err := service.DownloadAndSaveFromMagnet(ctx, mangetLink)
	require.NoError(t, err)

	err = service.restartDownloadForAllTorrents(ctx)
	require.NoError(t, err)
}

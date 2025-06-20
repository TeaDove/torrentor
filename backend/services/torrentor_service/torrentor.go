package torrentor_service

import (
	"context"
	"time"
	"torrentor/backend/schemas"
	"torrentor/backend/settings"

	"github.com/teadove/teasutils/utils/conv_utils"

	"github.com/teadove/teasutils/utils/must_utils"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) restartDownloadFromMagnet(
	ctx context.Context,
	magnetLink string,
) (*schemas.TorrentEntity, error) {
	// TODO move to settings
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	torrentObj, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, magnetLink)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download magnetLink")
	}

	torrentEnt, err := r.GetTorrentByInfoHash(ctx, torrentObj.InfoHash())
	if err != nil {
		return nil, errors.Wrap(err, "error making torrent object")
	}

	go must_utils.DoOrLogWithStacktrace(
		func(ctx context.Context) error { return r.onFileComplete(ctx, torrentEnt, time.Second*10) },
		"failed to run on torrent complete",
	)(ctx)

	return torrentEnt, nil
}

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	*schemas.TorrentEntity,
	<-chan torrent.TorrentStats,
	error,
) {
	torrentEnt, err := r.restartDownloadFromMagnet(ctx, magnetLink)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to download magnetLink")
	}

	zerolog.Ctx(ctx).
		Info().
		Object("torrent", torrentEnt).
		Msg("torrent.saved")

	return torrentEnt, r.torrentSupplier.ExportStats(ctx, torrentEnt.Obj), nil
}

type ServiceStats struct {
	StartedAt     time.Time
	TorrentsCount int
	FilesCount    int
	TotalSize     string
}

func (r *Service) makeTorrentStats(ctx context.Context) (ServiceStats, error) {
	torrents, err := r.GetAllTorrents(ctx)
	if err != nil {
		return ServiceStats{}, errors.Wrap(err, "failed to get torrents")
	}

	stats := ServiceStats{
		StartedAt:     settings.Settings.StartedAt,
		TorrentsCount: len(torrents),
	}
	var size conv_utils.Byte
	for _, torrentMeta := range torrents {
		size += torrentMeta.Size
		stats.FilesCount += len(torrentMeta.FilePathMap)
	}

	stats.TotalSize = size.String()

	return stats, nil
}

func (r *Service) Stats(ctx context.Context) (ServiceStats, <-chan torrent.ClientStats, error) {
	stats, err := r.makeTorrentStats(ctx)
	if err != nil {
		return ServiceStats{}, nil, errors.Wrap(err, "failed to make torrent stats")
	}

	return stats, r.torrentSupplier.Stats(ctx, time.Minute), nil
}

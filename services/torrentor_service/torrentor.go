package torrentor_service

import (
	"context"
	"time"
	"torrentor/schemas"

	"github.com/teadove/teasutils/utils/must_utils"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/settings_utils"
)

func (r *Service) restartDownload(
	ctx context.Context,
	magnetLink string,
) (*schemas.TorrentEntityPop, error) {
	torrentObj, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, magnetLink)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download magnetLink")
	}

	torrentEnt := r.makeTorrentMeta(torrentObj, magnetLink)

	torrentEnt, err = r.torrentRepository.TorrentUpsert(ctx, torrentEnt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save torrent")
	}

	torrentEntWithObj := &schemas.TorrentEntityPop{TorrentEntity: *torrentEnt, Obj: torrentObj}

	go must_utils.DoOrLogWithStacktrace(
		func(ctx context.Context) error { return r.onFileComplete(ctx, torrentEntWithObj, time.Second*10) },
		"failed to run on torrent complete",
	)(ctx)

	return torrentEntWithObj, nil
}

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	*schemas.TorrentEntity,
	<-chan torrent.TorrentStats,
	error,
) {
	torrentEnt, err := r.restartDownload(ctx, magnetLink)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to download magnetLink")
	}

	zerolog.Ctx(ctx).
		Info().
		Dict("torrent", torrentEnt.ZerologDict()).
		Msg("torrent.saved")

	return &torrentEnt.TorrentEntity, r.torrentSupplier.ExportStats(ctx, torrentEnt.Obj), nil
}

type ServiceStats struct {
	StartedAt     time.Time
	TorrentsCount int
	FilesCount    int
	TotalSize     uint64
}

func (r *Service) Stats(ctx context.Context) (ServiceStats, <-chan torrent.ClientStats, error) {
	torrents, err := r.torrentRepository.TorrentGetAll(ctx)
	if err != nil {
		return ServiceStats{}, nil, errors.Wrap(err, "failed to get torrents")
	}

	stats := ServiceStats{
		StartedAt:     settings_utils.BaseSettings.StartedAt,
		TorrentsCount: len(torrents),
	}
	for _, torrentMeta := range torrents {
		stats.TotalSize += torrentMeta.Root.Size
		stats.FilesCount += len(torrentMeta.Files)
	}

	return stats, r.torrentSupplier.Stats(ctx, time.Minute), nil
}

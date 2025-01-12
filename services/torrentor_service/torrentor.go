package torrentor_service

import (
	"context"
	"time"
	"torrentor/repositories/torrent_repository"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/settings_utils"
	"github.com/tidwall/buntdb"
)

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	torrent_repository.Torrent,
	<-chan torrent.TorrentStats,
	error,
) {
	torrentObj, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, magnetLink)
	if err != nil {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to download magnetLink")
	}

	torrentMeta, err := r.torrentRepository.TorrentGetByHash(ctx, torrentObj.InfoHash().String())
	if err == nil {
		zerolog.Ctx(ctx).
			Info().
			Dict("torrent", torrentMeta.ZerologDict()).
			Msg("torrent.already.exists")

		return torrentMeta, r.torrentSupplier.ExportStats(ctx, torrentObj), nil
	}

	if !errors.Is(err, buntdb.ErrNotFound) {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to get already created torrent")
	}

	torrentMeta = r.makeTorrentMeta(torrentObj, magnetLink)

	err = r.torrentRepository.TorrentSet(ctx, &torrentMeta)
	if err != nil {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to save torrent")
	}

	zerolog.Ctx(ctx).
		Info().
		Dict("torrent", torrentMeta.ZerologDict()).
		Msg("torrent.saved")

	return torrentMeta, r.torrentSupplier.ExportStats(ctx, torrentObj), nil
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

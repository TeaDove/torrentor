package torrentor_service

import (
	"context"
	"github.com/teadove/teasutils/utils/conv_utils"
	"gorm.io/gorm"
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

	torrentEnt, err := makeTorrentMeta(torrentObj)
	if err != nil {
		return nil, errors.Wrap(err, "error making torrent object")
	}

	err = r.torrentRepository.TorrentInsert(ctx, &torrentEnt)
	if err != nil {
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.Wrap(err, "failed to save torrent")
		}

		torrentEnt, err = r.torrentRepository.TorrentGetByHash(ctx, torrentObj.InfoHash().String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get torrent")
		}
	}

	torrentEntWithObj := &schemas.TorrentEntityPop{TorrentEntity: torrentEnt, Obj: torrentObj}

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
	TotalSize     conv_utils.Byte
}

func (r *Service) Stats(ctx context.Context) (ServiceStats, <-chan torrent.ClientStats, error) {
	torrents, err := r.GetAllTorrents(ctx)
	if err != nil {
		return ServiceStats{}, nil, errors.Wrap(err, "failed to get torrents")
	}

	stats := ServiceStats{
		StartedAt:     settings_utils.BaseSettings.StartedAt,
		TorrentsCount: len(torrents),
	}
	for _, torrentMeta := range torrents {
		stats.TotalSize += torrentMeta.Size
		stats.FilesCount += len(torrentMeta.Files)
	}

	return stats, r.torrentSupplier.Stats(ctx, time.Minute), nil
}

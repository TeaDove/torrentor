package torrentor_service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) restartDownloadForAllTorrents(ctx context.Context) error {
	torrents, err := r.listCreatedTorrents(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get all torrents")
	}

	stats, err := r.makeTorrentStats(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to make torrent stats")
	}

	zerolog.Ctx(ctx).
		Info().
		Interface("stats", stats).
		Msg("torrent.stats")

	for _, torrent := range torrents {
		zerolog.Ctx(ctx).
			Debug().
			Object("torrent", torrent).
			Msg("restarting.torrent")

		_, err = r.restartDownloadFromMagnet(ctx, torrent.Meta.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}
	}

	return nil
}

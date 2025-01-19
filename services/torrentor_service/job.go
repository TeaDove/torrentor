package torrentor_service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) restartDownloadForAllTorrents(ctx context.Context) error {
	torrents, err := r.GetAllTorrents(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get all torrents")
	}

	for _, torrent := range torrents {
		zerolog.Ctx(ctx).
			Debug().
			Dict("torrent", torrent.ZerologDict()).
			Msg("restarting.torrent")

		_, err = r.restartDownload(ctx, torrent.Meta.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}
	}

	return nil
}

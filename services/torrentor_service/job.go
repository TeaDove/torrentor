package torrentor_service

import (
	"context"
	"github.com/rs/zerolog"

	"github.com/pkg/errors"
)

func (r *Service) restartDownloadForAllTorrents(ctx context.Context) error {
	torrents, err := r.torrentRepository.TorrentGetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get torrents")
	}

	for _, torrent := range torrents {
		zerolog.Ctx(ctx).
			Debug().
			Dict("torrent", torrent.ZerologDict()).
			Msg("restarting.torrent")
		// TODO add restart from torrent without buntdb data
		_, err = r.restartDownload(ctx, torrent.Meta.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}
	}

	return nil
}

package torrentor_service

import (
	"context"
	"github.com/pkg/errors"
)

func (r *Service) RestartDownload(ctx context.Context) error {
	torrents, err := r.torrentRepository.TorrentGetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get torrents")
	}

	for _, torrent := range torrents {
		_, err = r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, torrent.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}
	}

	return nil
}

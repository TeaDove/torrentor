package torrentor_service

import (
	"context"
	"github.com/pkg/errors"
)

func (r *Service) restartDownload(ctx context.Context) error {
	torrents, err := r.torrentRepository.TorrentGetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get torrents")
	}

	for _, torrent := range torrents {
		_, err = r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, torrent.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}

		r.addOnTorrentCompleteCallback(ctx, &torrent, r.markCompleted)
	}

	return nil
}

func (r *Service) convertFormatsOnComplete(ctx context.Context) error {
	torrents, err := r.torrentRepository.TorrentGetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get torrents")
	}

	for _, torrent := range torrents {
		if torrent.Completed {
			continue
		}

		torrentSup, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, torrent.Magnet)
		if err != nil {
			return errors.Wrap(err, "unable to add magnet")
		}

		for _, file := range torrentSup.Files() {
			file.State()
		}
	}

	return nil
}

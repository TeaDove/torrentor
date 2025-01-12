package torrentor_service

import (
	"context"
	"github.com/rs/zerolog"
	"torrentor/repositories/torrent_repository"
)

func (r *Service) markCompleted(ctx context.Context, torrent *torrent_repository.Torrent) error {
	return r.torrentRepository.TorrentMarkComplete(ctx, torrent.ID)
}

func (r *Service) addOnTorrentCompleteCallback(
	ctx context.Context,
	torrent *torrent_repository.Torrent,
	callback func(context.Context, *torrent_repository.Torrent) error,
) {
	go func() {
		torrentSup, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, torrent.Magnet)
		if err != nil {
			zerolog.Ctx(ctx).
				Error().
				Stack().Err(err).
				Dict("torrent", torrent.ZerologDict()).
				Msg("torrent.callback.failed")
		}

		<-torrentSup.Complete().On()
		err = callback(ctx, torrent)
		if err != nil {
			zerolog.Ctx(ctx).
				Error().
				Stack().Err(err).
				Dict("torrent", torrent.ZerologDict()).
				Msg("torrent.callback.failed")
		}
	}()
}

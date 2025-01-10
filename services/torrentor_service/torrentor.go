package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/buntdb"
	"time"
	"torrentor/repositories/torrent_repository"
)

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	torrent_repository.Torrent,
	<-chan torrent.TorrentStats,
	error,
) {
	createdAt := time.Now().UTC()
	torrentObj, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, magnetLink)
	if err != nil {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to download magnetLink")
	}

	torrentMeta, err := r.torrentRepository.TorrentGetByHash(ctx, torrentObj.InfoHash().String())
	if err == nil {
		zerolog.Ctx(ctx).
			Info().
			Interface("torrent", &torrentMeta).
			Msg("torrent.already.exists")

		return torrentMeta, r.torrentSupplier.ExportStats(ctx, torrentObj), nil
	}

	if !errors.Is(err, buntdb.ErrNotFound) {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to get already created torrent")
	}

	id := uuid.New()
	torrentMeta = torrent_repository.Torrent{
		Id:          id,
		CreatedAt:   createdAt,
		Name:        torrentObj.Name(),
		Pieces:      uint64(torrentObj.NumPieces()),
		PieceLength: uint64(torrentObj.Info().PieceLength),
		InfoHash:    torrentObj.InfoHash().String(),
		Magnet:      magnetLink,
	}
	root := r.makeFile(torrentObj)

	torrentMeta.Root = root

	err = r.torrentRepository.TorrentSet(ctx, &torrentMeta)
	if err != nil {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to save torrent")
	}

	zerolog.Ctx(ctx).
		Info().
		Interface("torrent", &torrentMeta).
		Msg("torrent.saved")

	return torrentMeta, r.torrentSupplier.ExportStats(ctx, torrentObj), nil
}

func (r *Service) Stats(ctx context.Context) <-chan torrent.ClientStats {
	// TODO move to settings
	return r.torrentSupplier.Stats(ctx, time.Minute)
}

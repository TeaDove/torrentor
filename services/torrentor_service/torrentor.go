package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
	"torrentor/repositories/torrent_repository"
)

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	torrent_repository.Torrent,
	<-chan torrent.TorrentStats,
	error,
) {
	// TODO check for torrent existence
	createdAt := time.Now().UTC()
	torrentObj, err := r.torrentSupplier.AddMagnetAndGetInfoAndStartDownload(ctx, magnetLink)
	if err != nil {
		return torrent_repository.Torrent{}, nil, errors.Wrap(err, "failed to download magnetLink")
	}

	id := uuid.New()
	torrentMeta := torrent_repository.Torrent{
		Id:        id,
		CreatedAt: createdAt,
		Name:      torrentObj.Name(),
		Pieces:    uint64(torrentObj.NumPieces()),
		InfoHash:  torrentObj.InfoHash().HexString(),
		Magnet:    magnetLink,
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

func (r *Service) Stats() <-chan torrent.ClientStats {
	return r.torrentSupplier.Stats(time.Minute)
}

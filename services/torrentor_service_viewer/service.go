package torrentor_service_viewer

import (
	"context"
	"github.com/google/uuid"
	"torrentor/repositories/torrent_repository"
)

type Service struct {
	torrentRepository *torrent_repository.Repository
}

func NewService(_ context.Context, torrentRepository *torrent_repository.Repository) (*Service, error) {
	return &Service{torrentRepository: torrentRepository}, nil
}

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (torrent_repository.Torrent, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

package torrentor_service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"torrentor/repositories/torrent_repository"
)

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (torrent_repository.Torrent, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

func (r *Service) GetFileContentByID(ctx context.Context, id uuid.UUID) (*os.File, error) {
	filePath, err := r.torrentRepository.FileGetPath(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error getting file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}

	return file, nil
}

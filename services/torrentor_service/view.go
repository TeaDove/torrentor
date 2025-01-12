package torrentor_service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"path"
	"torrentor/repositories/torrent_repository"
	"torrentor/settings"
)

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (torrent_repository.Torrent, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

func (r *Service) GetFile(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (torrent_repository.FileWithContent, error) {
	torrent, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return torrent_repository.FileWithContent{}, errors.Wrap(err, "error getting torrent")
	}

	file, err := os.Open(path.Join(torrent.Location(settings.Settings.Torrent.DataDir), filePath))
	if err != nil {
		return torrent_repository.FileWithContent{}, errors.Wrap(err, "error opening file")
	}

	if _, ok := torrent.Files[filePath]; !ok {
		return torrent_repository.FileWithContent{}, errors.New("file not found")
	}

	return torrent_repository.FileWithContent{File: torrent.Files[filePath], OSFile: file}, nil
}

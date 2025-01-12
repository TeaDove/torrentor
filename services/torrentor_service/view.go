package torrentor_service

import (
	"context"
	"os"
	"path"
	"torrentor/repositories/torrent_repository"
	"torrentor/settings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Service) GetTorrentMetadataByID(ctx context.Context, id uuid.UUID) (torrent_repository.Torrent, error) {
	return r.torrentRepository.TorrentGetById(ctx, id)
}

func (r *Service) GetFileWithContent(
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

	fileMeta, ok := torrent.Files[filePath]
	if !ok {
		return torrent_repository.FileWithContent{}, errors.New("file not found")
	}

	return torrent_repository.FileWithContent{File: fileMeta, OSFile: file}, nil
}

func (r *Service) GetFile(
	ctx context.Context,
	torrentID uuid.UUID,
	filePath string,
) (torrent_repository.File, error) {
	torrent, err := r.torrentRepository.TorrentGetById(ctx, torrentID)
	if err != nil {
		return torrent_repository.File{}, errors.Wrap(err, "error getting torrent")
	}

	file, ok := torrent.Files[filePath]
	if !ok {
		return torrent_repository.File{}, errors.New("file not found")
	}

	return file, nil
}

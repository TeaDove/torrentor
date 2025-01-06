package torrentor_service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"time"
	"torrentor/repositories/torrent_repository"
)

func addToRoot(root *torrent_repository.File, file *torrent_repository.File) {
	// TODO set it advected
	if root.Files == nil {
		root.Files = make(map[string]torrent_repository.File, 1)
	}
	root.Files[file.Name] = *file
}

func setSize(file *torrent_repository.File) uint64 {
	if !file.IsDir {
		return file.Size
	}

	var child torrent_repository.File
	for _, child = range file.Files {
		file.Size += setSize(&child)
	}

	return file.Size
}

func (r *Service) populateTorrent(ctx context.Context, path string) (torrent_repository.File, error) {
	//TODO make path
	root := torrent_repository.File{
		Id:    uuid.New(),
		Name:  path,
		IsDir: true,
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to walk torrent")
		}

		file := torrent_repository.File{
			Id:   uuid.New(),
			Name: info.Name(),
			// TODO set total size for dirs
			Size:  uint64(info.Size()),
			IsDir: info.IsDir(),
		}
		addToRoot(&root, &file)

		return nil
	})
	if err != nil {
		return root, errors.Wrap(err, "filepath walk")
	}
	setSize(&root)

	return root, nil
}

func (r *Service) DownloadAndSaveFromMagnet(ctx context.Context, magnetLink string) (
	torrent_repository.Torrent,
	error,
) {
	createdAt := time.Now().UTC()
	torrentPath, err := r.torrentSupplier.DownloadMagnet(ctx, magnetLink)
	if err != nil {
		return torrent_repository.Torrent{}, errors.Wrap(err, "failed to download magnetLink")
	}

	id := uuid.New()
	torrent := torrent_repository.Torrent{Id: id, CreatedAt: createdAt, Name: torrentPath}
	root, err := r.populateTorrent(ctx, fmt.Sprintf("./data/torrent/%s", torrentPath))
	if err != nil {
		return torrent_repository.Torrent{}, errors.Wrap(err, "failed to populate torrent")
	}
	torrent.Root = root

	err = r.torrentRepository.TorrentSet(ctx, &torrent)
	if err != nil {
		return torrent_repository.Torrent{}, errors.Wrap(err, "failed to save torrent")
	}

	zerolog.Ctx(ctx).
		Info().
		Interface("torrent", &torrent).
		Msg("torrent.saved")

	return torrent, nil
}

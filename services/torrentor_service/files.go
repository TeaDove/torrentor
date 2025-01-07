package torrentor_service

import (
	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"torrentor/repositories/torrent_repository"
)

func addToRoot(root *torrent_repository.File, file *torrent_repository.File) {
	// TODO проставлять путь адеватно
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

func (r *Service) makeFile(torrentSup *torrent.Torrent) torrent_repository.File {
	root := torrent_repository.File{
		Id:    uuid.New(),
		Name:  torrentSup.Name(),
		IsDir: true,
	}

	for _, torrentFile := range torrentSup.Files() {
		if torrentFile == nil {
			continue
		}

		file := torrent_repository.File{
			Id:    uuid.New(),
			Name:  torrentFile.Path(), // TODO set correct path
			Size:  uint64(torrentFile.Length()),
			IsDir: false, // set dirs afterward
		}

		addToRoot(&root, &file)
	}

	setSize(&root)
	return root
}

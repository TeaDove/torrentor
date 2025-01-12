package torrentor_service

import (
	"mime"
	"path/filepath"
	"time"
	"torrentor/repositories/torrent_repository"

	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func (r *Service) makeTorrentMeta(torrentSup *torrent.Torrent, magnetLink string) torrent_repository.Torrent {
	id := uuid.New()
	createdAt := time.Now().UTC()

	torrentMeta := torrent_repository.Torrent{
		Id:          id,
		CreatedAt:   createdAt,
		Name:        torrentSup.Name(),
		Pieces:      uint64(torrentSup.NumPieces()),
		PieceLength: uint64(torrentSup.Info().PieceLength),
		InfoHash:    torrentSup.InfoHash().String(),
		Magnet:      magnetLink,
	}

	torrentMeta.Root = torrent_repository.File{
		Id:       uuid.New(),
		Name:     torrentSup.Name(),
		Path:     torrentSup.Name(),
		Mimetype: "",
		IsDir:    true,
	}
	torrentMeta.Files = make(map[string]torrent_repository.File, len(torrentSup.Files()))

	for _, torrentFile := range torrentSup.Files() {
		if torrentFile == nil {
			continue
		}

		file := torrent_repository.File{
			Id:       uuid.New(),
			Name:     filepath.Base(torrentFile.Path()),
			Path:     torrentFile.Path(),
			Mimetype: mime.TypeByExtension(filepath.Ext(torrentFile.Path())),
			Size:     uint64(torrentFile.Length()),
			SizeRepr: converters_utils.ToClosestByteAsString(torrentFile.Length(), 2),
			IsDir:    false,
		}

		torrentMeta.Files[file.Path] = file
		torrentMeta.Root.Size += file.Size
	}

	torrentMeta.Root.SizeRepr = converters_utils.ToClosestByteAsString(torrentMeta.Root.Size, 2)

	return torrentMeta
}

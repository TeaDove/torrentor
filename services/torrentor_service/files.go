package torrentor_service

import (
	"github.com/teadove/teasutils/utils/conv_utils"
	"mime"
	"path/filepath"
	"time"
	"torrentor/schemas"

	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
)

func makeFileEnt(path string, size uint64, completed bool) schemas.FileEntity {
	return schemas.FileEntity{
		Id:        uuid.New(),
		Name:      filepath.Base(path),
		Path:      schemas.TrimFirstDir(path),
		Mimetype:  mime.TypeByExtension(filepath.Ext(path)),
		Size:      conv_utils.Byte(size),
		Completed: completed,
	}
}

func (r *Service) makeTorrentMeta(torrentObj *torrent.Torrent, magnetLink string) *schemas.TorrentEntity {
	id := uuid.New()
	createdAt := time.Now().UTC()

	torrentMeta := schemas.TorrentEntity{
		ID:        id,
		CreatedAt: createdAt,
		Name:      torrentObj.Name(),
		Meta: schemas.Meta{
			Pieces:      uint64(torrentObj.NumPieces()),
			PieceLength: uint64(torrentObj.Info().PieceLength),
			Magnet:      magnetLink,
		},
		InfoHash: torrentObj.InfoHash().String(),
	}

	torrentMeta.Files = make(map[string]schemas.FileEntity, len(torrentObj.Files()))

	for _, torrentFile := range torrentObj.Files() {
		if torrentFile == nil {
			continue
		}

		file := makeFileEnt(
			torrentFile.Path(),
			uint64(torrentFile.Length()),
			false,
		)

		torrentMeta.AppendFile(file)
	}

	return &torrentMeta
}

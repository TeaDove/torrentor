package torrentor_service

import (
	"mime"
	"path/filepath"
	"time"
	"torrentor/schemas"

	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func makeFileEnt(path string, size uint64, completed bool) schemas.FileEntity {
	return schemas.FileEntity{
		Id:        uuid.New(),
		Name:      filepath.Base(path),
		Path:      path,
		Mimetype:  mime.TypeByExtension(filepath.Ext(path)),
		Size:      size,
		SizeRepr:  converters_utils.ToClosestByteAsString(size, 2),
		IsDir:     false,
		Completed: completed,
	}
}

func (r *Service) makeTorrentMeta(torrentObj *torrent.Torrent, magnetLink string) *schemas.TorrentEntity {
	id := uuid.New()
	createdAt := time.Now().UTC()

	torrentMeta := schemas.TorrentEntity{
		ID:      id,
		AddedAt: createdAt,
		Name:    torrentObj.Name(),
		Meta: schemas.Meta{
			Pieces:      uint64(torrentObj.NumPieces()),
			PieceLength: uint64(torrentObj.Info().PieceLength),
			Magnet:      magnetLink,
		},
		InfoHash: torrentObj.InfoHash().String(),
	}

	torrentMeta.Root = schemas.FileEntity{
		Id:       uuid.New(),
		Name:     torrentObj.Name(),
		Path:     torrentObj.Name(),
		Mimetype: "",
		IsDir:    true,
	}
	torrentMeta.Files = make(map[string]schemas.FileEntity, len(torrentObj.Files()))

	for _, torrentFile := range torrentObj.Files() {
		if torrentFile == nil {
			continue
		}

		file := makeFileEnt(torrentFile.Path(), uint64(torrentFile.Length()), false)

		torrentMeta.AppendFile(file)
		torrentMeta.Root.Size += file.Size
	}

	torrentMeta.Root.SizeRepr = converters_utils.ToClosestByteAsString(torrentMeta.Root.Size, 2)

	return &torrentMeta
}

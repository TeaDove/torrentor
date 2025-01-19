package torrentor_service

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/conv_utils"
	"mime"
	"path/filepath"
	"time"
	"torrentor/schemas"

	"github.com/anacrolix/torrent"
)

func makeFileEnt(path string, size uint64, completed bool) schemas.FileEntity {
	return schemas.FileEntity{
		Name:      filepath.Base(path),
		Path:      schemas.TrimFirstDir(path),
		Mimetype:  mime.TypeByExtension(filepath.Ext(path)),
		Size:      conv_utils.Byte(size),
		Completed: completed,
	}
}

func makeTorrentMeta(torrentObj *torrent.Torrent) (schemas.TorrentEntity, error) {
	createdAt := time.Now().UTC()

	metainfo := torrentObj.Metainfo()
	magnet, err := metainfo.MagnetV2()
	if err != nil {
		return schemas.TorrentEntity{}, errors.Wrap(err, "error getting magnet v2")
	}

	torrentMeta := schemas.TorrentEntity{
		CreatedAt: createdAt,
		Name:      torrentObj.Name(),
		Meta: schemas.Meta{
			Pieces:      uint64(torrentObj.NumPieces()),
			PieceLength: uint64(torrentObj.Info().PieceLength),
			Magnet:      magnet.String(),
		},
		InfoHash: torrentObj.InfoHash(),
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

	return torrentMeta, nil
}

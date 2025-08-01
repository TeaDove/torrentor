package torrentor_service

import (
	"context"
	"mime"
	"path/filepath"
	"time"
	"torrentor/backend/schemas"
	"torrentor/backend/utils/hash"

	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/conv_utils"

	"github.com/anacrolix/torrent"
)

func (r *Service) makeTorrentMeta(ctx context.Context, torrentObj *torrent.Torrent) (*schemas.TorrentEntity, error) {
	createdAt := time.Now().UTC()

	metainfo := torrentObj.Metainfo()

	magnet, err := metainfo.MagnetV2()
	if err != nil {
		return nil, errors.Wrap(err, "error getting magnet v2")
	}

	torrentMeta := schemas.TorrentEntity{
		CreatedAt: createdAt,
		Name:      torrentObj.Name(),
		Meta: schemas.Meta{
			Pieces:      uint64(torrentObj.NumPieces()),
			PieceLength: uint64(torrentObj.Info().PieceLength),
			Magnet:      magnet.String(),
		},
		InfoHash:       torrentObj.InfoHash(),
		TorrentDataDir: r.torrentDataDir,
		UnpackDataDir:  r.unpackDataDir,
		Obj:            torrentObj,
	}

	torrentMeta.FilePathMap = make(map[string]*schemas.FileEntity, len(torrentObj.Files()))
	torrentMeta.FileHashMap = make(map[string]*schemas.FileEntity, len(torrentObj.Files()))

	for _, torrentFile := range torrentObj.Files() {
		if torrentFile == nil {
			continue
		}

		path := schemas.TrimFirstDir(torrentFile.Path())
		file := &schemas.FileEntity{
			Name:      filepath.Base(torrentFile.Path()),
			Path:      path,
			PathHash:  hash.Sha1Base64Hash(path),
			Mimetype:  mime.TypeByExtension(filepath.Ext(torrentFile.Path())),
			Size:      conv_utils.Byte(torrentFile.Length()),
			Completed: false,
			Obj:       torrentFile,
			Torrent:   &torrentMeta,
		}

		file.Meta, err = r.ffmpegService.ExportMetadata(ctx, file.Location())
		if err != nil {
			return nil, errors.Wrap(err, "error getting metadata")
		}

		torrentMeta.AppendFile(file)
	}

	return &torrentMeta, nil
}

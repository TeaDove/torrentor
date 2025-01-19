package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"os"
	"torrentor/schemas"
)

func (r *Service) GetFileWithContent(
	ctx context.Context,
	torrentInfoHash metainfo.Hash,
	filePath string,
) (schemas.FileWithContent, error) {
	torrent, err := r.GetTorrentByInfoHash(ctx, torrentInfoHash)
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "error getting torrent")
	}

	file, err := os.Open(torrent.FileLocation(r.torrentDataDir, filePath))
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "error opening file")
	}

	fileMeta, ok := torrent.Files[filePath]
	if !ok {
		return schemas.FileWithContent{}, errors.New("file not found")
	}

	return schemas.FileWithContent{FileEntity: fileMeta, OSFile: file}, nil
}

func (r *Service) GetFile(
	ctx context.Context,
	torrentInfoHash metainfo.Hash,
	filePath string,
) (*schemas.FileEntity, error) {
	torrentEnt, err := r.GetTorrentByInfoHash(ctx, torrentInfoHash)
	if err != nil {
		return nil, errors.Wrap(err, "error getting torrent")
	}

	file, ok := torrentEnt.Files[filePath]
	if !ok {
		return nil, errors.New("file not found")
	}

	return file, nil
}

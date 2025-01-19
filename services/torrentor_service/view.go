package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"torrentor/schemas"
)

func (r *Service) GetFileWithContent(
	ctx context.Context,
	torrentInfoHash metainfo.Hash,
	filePath string,
) (schemas.FileWithContent, error) {
	fileEnt, err := r.GetFileByInfoHashAndPath(ctx, torrentInfoHash, filePath)
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "failed to get torrent info")
	}

	file, err := os.Open(fileEnt.Location())
	if err != nil {
		return schemas.FileWithContent{}, errors.Wrap(err, "error opening file")
	}

	return schemas.FileWithContent{FileEntity: fileEnt, OSFile: file}, nil
}

func (r *Service) GetHLS(
	ctx context.Context,
	torrentInfoHash metainfo.Hash,
	fileHash string,
) (string, error) {
	fileEnt, err := r.GetFileByInfoHashAndHash(ctx, torrentInfoHash, fileHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get torrent info")
	}

	metadata, err := r.ffmpegService.ExportMetadata(ctx, fileEnt.Location())
	if err != nil {
		return "", errors.Wrap(err, "error exporting metadata")
	}

	unpackFilesDir := filepath.Dir(fileEnt.LocationInUnpack())
	location := path.Join(unpackFilesDir, makeFilenameWithTags(
		"hls",
		"/",
		metadata.Streams[0].Tags.Title,
		metadata.Streams[0].Tags.Language,
	))

	return location, nil
}

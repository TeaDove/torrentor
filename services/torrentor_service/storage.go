package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
	"torrentor/schemas"
)

func (r *Service) GetFileByInfoHashAndPath(ctx context.Context, infoHash metainfo.Hash, filepath string) (*schemas.FileEntity, error) {
	torrentEnt, err := r.GetTorrentByInfoHash(ctx, infoHash)
	if err != nil {
		return nil, errors.Wrap(err, "error getting torrent by hash")
	}

	file, ok := torrentEnt.FilePathMap[filepath]
	if !ok {
		return nil, errors.New("file not found")
	}

	return file, nil
}

func (r *Service) UnpackIfNeeded(ctx context.Context, fileEnt *schemas.FileEntity) error {
	t0 := time.Now()
	if fileEnt.Mimetype == schemas.MatroskaMimeType {
		err := r.unpackMatroska(ctx, fileEnt)
		if err != nil {
			return errors.Wrap(err, "failed to unpack matroska file")
		}
	}

	zerolog.Ctx(ctx).
		Info().
		Dict("file", fileEnt.ZerologDict()).
		Str("elapsed", time.Since(t0).String()).
		Msg("unpacked")

	return nil
}

func (r *Service) GetFileByInfoHashAndHash(ctx context.Context, infoHash metainfo.Hash, filehash string) (*schemas.FileEntity, error) {
	torrentEnt, err := r.GetTorrentByInfoHash(ctx, infoHash)
	if err != nil {
		return nil, errors.Wrap(err, "error getting torrent by hash")
	}

	file, ok := torrentEnt.FileHashMap[filehash]
	if !ok {
		return nil, errors.New("file not found")
	}

	return file, nil
}

func (r *Service) GetTorrentByInfoHash(ctx context.Context, infoHash metainfo.Hash) (*schemas.TorrentEntity, error) {
	r.hashToTorrentMu.RLock()
	torrentEnt, ok := r.hashToTorrent[infoHash]
	r.hashToTorrentMu.RUnlock()

	if ok {
		return torrentEnt, nil
	}

	torrentObj, err := r.torrentSupplier.GetTorrentByInfoHash(ctx, infoHash)
	if err != nil {
		return &schemas.TorrentEntity{}, errors.Wrap(err, "error getting torrent obj")
	}

	torrentMeta, err := r.makeTorrentMeta(ctx, torrentObj)
	if err != nil {
		return &schemas.TorrentEntity{}, errors.Wrap(err, "error getting torrent metadata")
	}

	r.hashToTorrentMu.Lock()
	r.hashToTorrent[infoHash] = torrentMeta
	r.hashToTorrentMu.Unlock()

	return torrentMeta, nil
}

func (r *Service) GetAllTorrents(ctx context.Context) ([]*schemas.TorrentEntity, error) {
	torrentsDir, err := os.ReadDir(r.torrentDataDir)
	if err != nil {
		return nil, errors.Wrap(err, "error reading torrent dir")
	}

	var (
		torrents   = make([]*schemas.TorrentEntity, 0, 5)
		torrentEnt *schemas.TorrentEntity
	)

	for _, file := range torrentsDir {
		if !file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		hash := metainfo.Hash{}
		err = hash.FromHexString(file.Name())
		if err != nil {
			zerolog.Ctx(ctx).
				Error().
				Stack().Err(err).
				Str("infoHash", file.Name()).
				Msg("failed.to.parse.info.hash")
		}

		torrentEnt, err = r.GetTorrentByInfoHash(ctx, hash)
		if err != nil {
			zerolog.Ctx(ctx).
				Error().
				Stack().Err(err).
				Str("infoHash", file.Name()).
				Msg("failed.to.get.torrent.by.info.hash")
		}

		torrents = append(torrents, torrentEnt)
	}

	return torrents, nil
}

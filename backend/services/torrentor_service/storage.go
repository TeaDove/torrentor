package torrentor_service

import (
	"context"
	"os"
	"strings"
	"time"
	"torrentor/backend/schemas"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Service) GetFileByInfoHashAndPath(
	ctx context.Context,
	infoHash metainfo.Hash,
	filepath string,
) (*schemas.FileEntity, error) {
	torrentEnt, err := r.GetOrCreateTorrentByInfoHash(ctx, infoHash)
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
		Object("file", fileEnt).
		Str("elapsed", time.Since(t0).String()).
		Msg("unpacked")

	return nil
}

func (r *Service) GetTorrentByInfoHash(infoHash metainfo.Hash) (*schemas.TorrentEntity, bool) {
	r.hashToTorrentMu.RLock()
	torrentEnt, ok := r.hashToTorrent[infoHash]
	r.hashToTorrentMu.RUnlock()

	return torrentEnt, ok
}

func (r *Service) GetOrCreateTorrentByInfoHash(
	ctx context.Context,
	infoHash metainfo.Hash,
) (*schemas.TorrentEntity, error) {
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

func (r *Service) DeleteTorrentByInfoHash(infoHash metainfo.Hash) bool {
	r.hashToTorrentMu.RLock()
	torrentEnt, ok := r.hashToTorrent[infoHash]
	r.hashToTorrentMu.RUnlock()

	if !ok {
		return false
	}

	r.hashToTorrentMu.Lock()
	torrentEnt.Obj.Drop()
	delete(r.hashToTorrent, infoHash)

	_ = os.RemoveAll(torrentEnt.RootLocation())
	_ = os.RemoveAll(torrentEnt.UnpackLocation())
	r.hashToTorrentMu.Unlock()

	return true
}

func (r *Service) ListOpenTorrents(ctx context.Context) ([]*schemas.TorrentEntity, error) {
	torrentsDir, err := os.ReadDir(r.torrentDataDir)
	if err != nil {
		return nil, errors.Wrap(err, "error reading torrent dir")
	}

	var (
		torrents   = make([]*schemas.TorrentEntity, 0, 5)
		torrentEnt *schemas.TorrentEntity
		ok         bool
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

			continue
		}

		torrentEnt, ok = r.GetTorrentByInfoHash(hash)
		if !ok {
			zerolog.Ctx(ctx).
				Error().
				Str("infoHash", file.Name()).
				Msg("torrent.not.found")

			continue
		}

		torrents = append(torrents, torrentEnt)
	}

	return torrents, nil
}

func (r *Service) listCreatedTorrents(ctx context.Context) ([]*schemas.TorrentEntity, error) {
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
			zerolog.Ctx(ctx).Error().
				Stack().Err(err).
				Str("infoHash", file.Name()).
				Msg("failed.to.parse.info.hash")

			continue
		}

		torrentEnt, err = r.GetOrCreateTorrentByInfoHash(ctx, hash)
		if err != nil {
			zerolog.Ctx(ctx).Error().
				Stack().Err(err).
				Str("infoHash", file.Name()).
				Msg("torrent.not.found")

			continue
		}

		torrents = append(torrents, torrentEnt)
	}

	return torrents, nil
}

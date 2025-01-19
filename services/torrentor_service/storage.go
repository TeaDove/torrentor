package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"torrentor/schemas"
)

func (r *Service) GetTorrentByInfoHash(ctx context.Context, infoHash metainfo.Hash) (schemas.TorrentEntityPop, error) {
	torrentObj, err := r.torrentSupplier.GetTorrentByInfoHash(ctx, infoHash)
	if err != nil {
		return schemas.TorrentEntityPop{}, errors.Wrap(err, "error getting torrent obj")
	}

	torrentMeta, err := makeTorrentMeta(torrentObj)
	if err != nil {
		return schemas.TorrentEntityPop{}, errors.Wrap(err, "error getting torrent metadata")
	}

	torrentEnt := schemas.TorrentEntityPop{
		TorrentEntity: torrentMeta,
		Obj:           torrentObj,
	}
	// TODO add files from torrent folder

	return torrentEnt, nil
}

func (r *Service) GetAllTorrents(ctx context.Context) ([]schemas.TorrentEntityPop, error) {
	torrentsDir, err := os.ReadDir(r.torrentDataDir)
	if err != nil {
		return nil, errors.Wrap(err, "error reading torrent dir")
	}

	var (
		torrents   = make([]schemas.TorrentEntityPop, 0, 5)
		torrentEnt schemas.TorrentEntityPop
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

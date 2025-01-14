package torrent_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"torrentor/schemas"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/buntdb"
)

func makeInfoHashToTorrentKey(infoHash string) string {
	return fmt.Sprintf("torrent:%s", infoHash)
}

func (r *Repository) TorrentGetByHash(_ context.Context, infoHash string) (schemas.TorrentEntity, error) {
	val, err := r.db.Get(makeInfoHashToTorrentKey(infoHash))
	if err != nil {
		return schemas.TorrentEntity{}, errors.Wrap(err, "failed to get torrent by link")
	}

	var torrentEnt schemas.TorrentEntity

	err = json.Unmarshal([]byte(val), &torrentEnt)
	if err != nil {
		return schemas.TorrentEntity{}, errors.Wrap(err, "failed to unmarshal torrent by link")
	}

	return torrentEnt, nil
}

func (r *Repository) TorrentGetById(_ context.Context, id uuid.UUID) (schemas.TorrentEntity, error) {
	val, err := r.db.GetByIndex(torrentByIDIdx, fmt.Sprintf(`{"id":"%s"}`, id))
	if err != nil {
		return schemas.TorrentEntity{}, errors.Wrap(err, "failed to get torrent by id")
	}

	var torrentEnt schemas.TorrentEntity

	err = json.Unmarshal([]byte(val), &torrentEnt)
	if err != nil {
		return schemas.TorrentEntity{}, errors.Wrap(err, "failed to unmarshal torrent by link")
	}

	return torrentEnt, nil
}

func (r *Repository) torrentSet(tx *buntdb.Tx, torrent *schemas.TorrentEntity) error {
	val, err := json.Marshal(&torrent)
	if err != nil {
		return errors.Wrap(err, "failed to marshal torrent")
	}

	_, _, err = tx.Set(makeInfoHashToTorrentKey(torrent.InfoHash), string(val), nil)
	if err != nil {
		return errors.Wrap(err, "failed to set torrent")
	}

	return nil
}

func (r *Repository) TorrentUpsert(_ context.Context, torrent *schemas.TorrentEntity) (*schemas.TorrentEntity, error) {
	var err error

	err = r.db.Update(func(tx *buntdb.Tx) error {
		key := makeInfoHashToTorrentKey(torrent.InfoHash)

		var v string

		v, err = tx.Get(key)
		if errors.Is(err, buntdb.ErrNotFound) {
			return r.torrentSet(tx, torrent)
		}

		if err != nil {
			return errors.Wrap(err, "failed to get torrent")
		}

		var oldTorrent schemas.TorrentEntity

		err = json.Unmarshal([]byte(v), &oldTorrent)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal torrent")
		}

		torrent.ID = oldTorrent.ID

		return r.torrentSet(tx, torrent)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert torrent")
	}

	return torrent, nil
}

func (r *Repository) TorrentGetAll(ctx context.Context) ([]schemas.TorrentEntity, error) {
	var torrents []schemas.TorrentEntity

	err := r.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend(torrentByIDIdx, func(key, value string) bool {
			var torrent schemas.TorrentEntity

			err := json.Unmarshal([]byte(value), &torrent)
			if err != nil {
				zerolog.Ctx(ctx).
					Error().
					Stack().Err(err).
					Str("v", value).
					Msg("failed.to.unmarshal.torrent")

				return true
			}

			torrents = append(torrents, torrent)

			return true
		})

		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get torrents")
	}

	return torrents, nil
}

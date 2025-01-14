package torrent_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
	"torrentor/schemas"
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

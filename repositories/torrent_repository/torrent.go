package torrent_repository

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

func makeTorrentKey(id uuid.UUID) string {
	return fmt.Sprintf("torrent:%s", id)
}

func (r *Repository) TorrentGetById(ctx context.Context, id uuid.UUID) (Torrent, error) {
	var torrent Torrent

	err := r.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(makeTorrentKey(id))
		if err != nil {
			if errors.Is(err, buntdb.ErrNotFound) {
				return stderrors.Join(ErrNotFound, err)
			}

			return errors.Wrap(err, "failed to get torrent by link")
		}

		err = json.Unmarshal([]byte(val), &torrent)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal torrent by link")
		}

		return nil
	})
	if err != nil {
		return torrent, errors.Wrap(err, "error getting torrent")
	}

	return torrent, nil
}

func (r *Repository) TorrentSet(ctx context.Context, torrent *Torrent) error {
	err := r.db.Update(func(tx *buntdb.Tx) error {
		val, err := json.Marshal(torrent)
		if err != nil {
			return errors.Wrap(err, "failed to marshal torrent")
		}

		// TODO create ttl
		_, _, err = tx.Set(makeTorrentKey(torrent.Id), string(val), nil)
		if err != nil {
			return errors.Wrap(err, "failed to save torrent")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error updating torrent")
	}

	return nil
}

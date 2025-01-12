package torrent_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/buntdb"
)

func makeInfoHashToTorrentKey(infoHash string) string {
	return fmt.Sprintf("torrent:%s", infoHash)
}

func (r *Repository) TorrentGetByHash(_ context.Context, infoHash string) (Torrent, error) {
	val, err := r.db.Get(makeInfoHashToTorrentKey(infoHash))
	if err != nil {
		return Torrent{}, errors.Wrap(err, "failed to get torrent by link")
	}

	var torrent Torrent
	err = json.Unmarshal([]byte(val), &torrent)
	if err != nil {
		return Torrent{}, errors.Wrap(err, "failed to unmarshal torrent by link")
	}

	return torrent, nil
}

func (r *Repository) TorrentGetById(_ context.Context, id uuid.UUID) (Torrent, error) {
	var torrent Torrent

	val, err := r.db.GetByIndex(torrentByIDIdx, fmt.Sprintf(`{"id":"%s"}`, id))
	if err != nil {
		return Torrent{}, errors.Wrap(err, "failed to get torrent by id")
	}

	err = json.Unmarshal([]byte(val), &torrent)
	if err != nil {
		return Torrent{}, errors.Wrap(err, "failed to unmarshal torrent by link")
	}

	return torrent, nil
}

func (r *Repository) TorrentSet(_ context.Context, torrent *Torrent) error {
	val, err := json.Marshal(torrent)
	if err != nil {
		return errors.Wrap(err, "failed to marshal torrent")
	}

	// TODO create ttl
	_, _, err = r.db.Set(makeInfoHashToTorrentKey(torrent.InfoHash), string(val), nil)
	if err != nil {
		return errors.Wrap(err, "failed to save torrent")
	}

	return nil
}

func (r *Repository) TorrentGetAll(ctx context.Context) ([]Torrent, error) {
	var torrents []Torrent
	err := r.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend(torrentByIDIdx, func(key, value string) bool {
			var torrent Torrent
			err := json.Unmarshal([]byte(value), &torrent)
			if err != nil {
				zerolog.Ctx(ctx).Error().Stack().Err(err).Str("v", value).Msg("failed.to.unmarshal.torrent")
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

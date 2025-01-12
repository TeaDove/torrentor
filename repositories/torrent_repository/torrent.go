package torrent_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

	val, err := r.db.GetByIndex(torrentById, fmt.Sprintf(`{"id":"%s"}`, id))
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

package torrent_repository

import (
	"context"
	"torrentor/schemas"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Repository) TorrentGetByHash(_ context.Context, infoHash string) (schemas.TorrentEntity, error) {
	var torrent schemas.TorrentEntity
	err := r.db.Where("info_hash = ?", infoHash).Find(&torrent).Error
	if err != nil {
		return torrent, errors.Wrap(err, "error getting torrent")
	}

	return torrent, nil
}

func (r *Repository) TorrentGetById(_ context.Context, id uuid.UUID) (schemas.TorrentEntity, error) {
	var torrent schemas.TorrentEntity
	err := r.db.Where("id = ?", id).Find(&torrent).Error
	if err != nil {
		return torrent, errors.Wrap(err, "error getting torrent")
	}

	return torrent, nil
}

//func (r *Repository) torrentSet(tx *buntdb.Tx, torrent *schemas.TorrentEntity) error {
//	val, err := json.Marshal(&torrent)
//	if err != nil {
//		return errors.Wrap(err, "failed to marshal torrent")
//	}
//
//	_, _, err = tx.Set(makeInfoHashToTorrentKey(torrent.InfoHash), string(val), nil)
//	if err != nil {
//		return errors.Wrap(err, "failed to set torrent")
//	}
//
//	return nil
//}

func (r *Repository) TorrentSave(_ context.Context, torrent *schemas.TorrentEntity) (*schemas.TorrentEntity, error) {
	err := r.db.Save(torrent).Error
	if err != nil {
		return nil, errors.Wrap(err, "error updating torrent")
	}

	return torrent, nil
}

func (r *Repository) TorrentGetAll(ctx context.Context) ([]schemas.TorrentEntity, error) {
	torrents := make([]schemas.TorrentEntity, 0)

	err := r.db.Find(&torrents).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get torrents")
	}

	return torrents, nil
}

package torrent_repository

import (
	"context"
	"torrentor/schemas"

	"github.com/google/uuid"
)

func (r *Repository) TorrentGetByHash(ctx context.Context, infoHash string) (schemas.TorrentEntity, error) {
	var torrent schemas.TorrentEntity
	//err := r.db.
	//	WithContext(ctx).
	//	Where("info_hash = ?", infoHash).
	//	Find(&torrent).
	//	Error
	//if err != nil {
	//	return torrent, errors.Wrap(err, "error getting torrent")
	//}

	return torrent, nil
}

func (r *Repository) TorrentGetAll(ctx context.Context) ([]schemas.TorrentEntity, error) {
	torrents := make([]schemas.TorrentEntity, 0)

	//err := r.db.
	//	WithContext(ctx).
	//	Find(&torrents).
	//	Error
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to get torrents")
	//}

	return torrents, nil
}

func (r *Repository) TorrentInsert(ctx context.Context, torrent *schemas.TorrentEntity) error {
	//err := r.db.
	//	WithContext(ctx).
	//	Create(torrent).
	//	Error
	//if err != nil {
	//	return errors.Wrap(err, "error saving torrent")
	//}

	return nil
}

func (r *Repository) TorrentUpdate(ctx context.Context, id uuid.UUID, values map[string]any) error {
	//err := r.db.
	//	WithContext(ctx).
	//	Updates(values).
	//	Where("id = ?", id).
	//	Error
	//if err != nil {
	//	return errors.Wrap(err, "error saving torrent")
	//}

	return nil
}

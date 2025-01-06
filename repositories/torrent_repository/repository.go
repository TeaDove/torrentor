package torrent_repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

type Repository struct {
	db *buntdb.DB
}

func NewRepository(ctx context.Context, db *buntdb.DB) (*Repository, error) {
	// TODO move to settings
	// TODO get from repo
	return &Repository{db: db}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.db.Close()
}

func (r *Repository) Health(ctx context.Context) error {
	err := r.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(`SELECT_1`)
		if err != nil {
			if errors.Is(err, buntdb.ErrNotFound) {
				return nil
			}
			return errors.Wrap(err, `failed to get 1`)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to check health")
	}

	return nil
}

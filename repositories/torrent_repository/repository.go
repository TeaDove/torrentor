package torrent_repository

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"torrentor/schemas"

	"github.com/pkg/errors"
)

type Repository struct {
	db    *gorm.DB
	sqldb *sql.DB
}

const torrentByIDIdx = "torrent-by-id"

func NewRepository(_ context.Context, db *gorm.DB) (*Repository, error) {
	sqldb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get db sql")
	}

	r := &Repository{db: db, sqldb: sqldb}

	err = r.db.AutoMigrate(&schemas.TorrentEntity{}, &schemas.FileEntity{})
	if err != nil {
		return nil, errors.Wrap(err, "auto migrate sqlite")
	}

	return r, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.sqldb.Close()
}

func (r *Repository) Health(ctx context.Context) error {
	return r.sqldb.PingContext(ctx)
}

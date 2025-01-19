package torrent_repository

import (
	"context"
)

type Repository struct {
	dataDir string
}

func NewRepository(_ context.Context) (*Repository, error) {
	return &Repository{}, nil
}

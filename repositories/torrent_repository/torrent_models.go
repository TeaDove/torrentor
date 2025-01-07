package torrent_repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Torrent struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Magnet    string    `json:"magnet"`

	Root        File   `json:"root,omitempty"`
	Pieces      uint64 `json:"pieces,omitempty"`
	PieceLength uint64 `json:"piecesLength,omitempty"`
	InfoHash    string `json:"infoHash"`
}

type File struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Size uint64    `json:"size"`

	IsDir bool            `json:"isDir"`
	Files map[string]File `json:"files,omitempty"`
}

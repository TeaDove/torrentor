package torrent_repository

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/converters_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
	"time"
)

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

func (r *Torrent) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("id", r.Id.String()).
		Str("name", r.Name).
		Str("magnet", redact_utils.Trim(r.Magnet)).
		Float64("size_mb", converters_utils.ToMegaByte(r.Root.Size))
}

type File struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Size uint64    `json:"size"`

	IsDir bool            `json:"isDir"`
	Files map[string]File `json:"files,omitempty"`
}

func (r *File) FlatFiles() []File {
	files := []File{*r}
	for _, file := range r.Files {
		files = append(files, file.FlatFiles()...)
	}

	return files
}

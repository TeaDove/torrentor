package torrent_repository

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
	"maps"
	"path/filepath"
	"slices"
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
		Str("size", r.Root.SizeRepr)
}

type File struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Size     uint64    `json:"size"`
	SizeRepr string    `json:"sizeRepr"`

	IsDir bool            `json:"isDir"`
	Files map[string]File `json:"files,omitempty"`
}

func fileCompare(a File, b File) int {
	if a.IsDir != b.IsDir {
		if a.IsDir {
			return 1
		} else {
			return -1
		}
	}

	if a.Name == b.Name {
		return 0
	}

	aExt, bExt := filepath.Ext(a.Name), filepath.Ext(b.Name)
	if aExt == bExt {
		if a.Name > b.Name {
			return 1
		} else {
			return -1
		}
	}

	if aExt > bExt {
		return 1
	} else {
		return -1
	}
}

func (r *File) FlatFiles() []File {
	files := []File{*r}

	sortedFiles := slices.SortedFunc(maps.Values(r.Files), fileCompare)

	for _, file := range sortedFiles {
		files = append(files, file.FlatFiles()...)
	}

	return files
}

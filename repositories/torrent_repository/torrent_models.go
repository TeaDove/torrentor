package torrent_repository

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type Torrent struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Magnet    string    `json:"magnet"`

	Root  File            `json:"root,omitempty"`
	Files map[string]File `json:"files,omitempty"`

	Pieces      uint64 `json:"pieces,omitempty"`
	PieceLength uint64 `json:"piecesLength,omitempty"`
	InfoHash    string `json:"infoHash"`
	Completed   bool   `json:"completed"`
}

func (r *Torrent) MarshalJSON() ([]byte, error) {
	val, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal torrent")
	}

	return val, nil
}

func (r *Torrent) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal torrent")
	}

	return nil
}

func (r *Torrent) Location(dataDir string) string {
	return path.Join(dataDir, r.InfoHash)
}

func (r *Torrent) FileLocation(dataDir string, filePath string) string {
	return path.Join(r.Location(dataDir), filePath)
}

func (r *Torrent) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("id", r.ID.String()).
		Str("name", r.Name).
		Str("magnet", redact_utils.Trim(r.Magnet)).
		Str("size", r.Root.SizeRepr)
}

type File struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Mimetype string    `json:"mimetype,omitempty"`
	Size     uint64    `json:"size"`
	SizeRepr string    `json:"sizeRepr"`

	IsDir bool `json:"isDir"`
}

func (r *File) IsVideo() bool {
	return strings.Split(r.Mimetype, "/")[0] == "video"
}

type FileWithContent struct {
	File
	OSFile *os.File
}

func fileCompare(a File, b File) int {
	if a.IsDir != b.IsDir {
		if a.IsDir {
			return 1
		} else {
			return -1
		}
	}

	if a.Path == b.Path {
		return 0
	}

	aExt, bExt := filepath.Ext(a.Path), filepath.Ext(b.Path)
	if aExt == bExt {
		if a.Path > b.Path {
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

func (r *Torrent) FlatFiles() []File {
	return slices.SortedFunc(maps.Values(r.Files), fileCompare)
}

const (
	MatroskaMimeType = "video/x-matroska"
)

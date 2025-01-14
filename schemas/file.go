package schemas

import (
	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileEntity struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Path string    `json:"path"`

	Mimetype string `json:"mimetype,omitempty"`
	Size     uint64 `json:"size"`
	SizeRepr string `json:"sizeRepr"`

	IsDir     bool `json:"isDir"`
	Completed bool `json:"completed"`
}

func (r *FileEntityPop) BaseName() string {
	return filepath.Base(r.Name)
}

func (r *FileEntityPop) NameWithoutExt() string {
	return strings.TrimSuffix(r.BaseName(), filepath.Ext(r.Name))
}

func (r *FileEntity) IsVideo() bool {
	return strings.Split(r.Mimetype, "/")[0] == "video"
}

type FileEntityPop struct {
	FileEntity
	Obj     *torrent.File     `json:"-"`
	Torrent *TorrentEntityPop `json:"-"`
}

func (r *FileEntityPop) Location(dataDir string) string {
	return r.Torrent.FileLocation(dataDir, r.Path)
}

type FileWithContent struct {
	FileEntity
	OSFile *os.File `json:"-"`
}

func fileCompare(a FileEntity, b FileEntity) int {
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

func (r *TorrentEntity) FlatFiles() []FileEntity {
	return slices.SortedFunc(maps.Values(r.Files), fileCompare)
}

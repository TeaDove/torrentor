package schemas

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/conv_utils"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/anacrolix/torrent"
)

type FileEntity struct {
	Name string `json:"name"`
	Path string `json:"path"`

	Mimetype string          `json:"mimetype,omitempty"`
	Size     conv_utils.Byte `json:"size"`

	Completed bool `json:"completed"`
}

func TrimFirstDir(path string) string {
	fields := strings.Split(path, "/")
	if len(fields) < 2 {
		panic(errors.Errorf("invalid path: %s", path))
	}

	return filepath.Join(fields[1:]...)
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
	if a.Path > b.Path {
		return 1
	} else {
		return -1
	}

	//if a.IsDir != b.IsDir {
	//	if a.IsDir {
	//		return 1
	//	} else {
	//		return -1
	//	}
	//}

	//aFields, bFields := strings.Split(a.Path, "/"), strings.Split(b.Path, "/")
	//idx := 0
	//for {
	//	if aFields[idx] != bFields[idx] {}
	//	idx += 1
	//}

	//if a.Path == b.Path {
	//	return 0
	//}
	//
	//aExt, bExt := filepath.Ext(a.Path), filepath.Ext(b.Path)
	//if aExt == bExt {
	//	if a.Path > b.Path {
	//		return 1
	//	} else {
	//		return -1
	//	}
	//}
	//
	//if aExt > bExt {
	//	return 1
	//} else {
	//	return -1
	//}
}

func (r *TorrentEntity) FlatFiles() []FileEntity {
	return slices.SortedFunc(maps.Values(r.Files), fileCompare)
}

package schemas

import (
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"torrentor/backend/services/ffmpeg_service"

	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/conv_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
)

type FileEntity struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	PathHash string `json:"pathHash"`

	Mimetype string          `json:"mimetype,omitempty"`
	Size     conv_utils.Byte `json:"size"`

	Completed bool `json:"completed"`

	Meta    ffmpeg_service.Metadata `json:"meta,omitempty"`
	Obj     *torrent.File           `json:"-"`
	Torrent *TorrentEntity          `json:"-"`
}

func (r *FileEntity) MarshalZerologObject(e *zerolog.Event) {
	e.Str("hash", redact_utils.Trim(r.PathHash)).Str("name", r.Name).Str("size", r.Size.String())
}

func TrimFirstDir(path string) string {
	fields := strings.Split(path, "/")
	if len(fields) < 2 {
		panic(errors.Errorf("invalid path: %s", path))
	}

	return filepath.Join(fields[1:]...)
}

func (r *FileEntity) BaseName() string {
	return filepath.Base(r.Name)
}

func (r *FileEntity) NameWithoutExt() string {
	return strings.TrimSuffix(r.BaseName(), filepath.Ext(r.Name))
}

func (r *FileEntity) IsVideo() bool {
	return strings.Split(r.Mimetype, "/")[0] == "video"
}

func (r *FileEntity) Location() string {
	return path.Join(r.Torrent.RawLocation(), r.Path)
}

func (r *FileEntity) LocationInUnpack() string {
	return path.Join(r.Torrent.UnpackLocation(), r.PathHash)
}

func (r *FileEntity) LocationInUnpackAsStream(stream *ffmpeg_service.Stream, suffix string) string {
	return path.Join(r.Torrent.UnpackLocation(), r.PathHash, stream.StreamFile(suffix))
}

type FileWithContent struct {
	*FileEntity
	OSFile *os.File `json:"-"`
}

func fileCompare(a *FileEntity, b *FileEntity) int {
	if a.Path > b.Path {
		return 1
	}

	return -1
	//	if a.IsDir != b.IsDir {
	//		if a.IsDir {
	//			return 1
	//		} else {
	//			return -1
	//		}
	//	}
	//
	// aFields, bFields := strings.Split(a.Path, "/"), strings.Split(b.Path, "/")
	// idx := 0
	//
	//	for {
	//		if aFields[idx] != bFields[idx] {}
	//		idx += 1
	//	}
	//
	//	if a.Path == b.Path {
	//		return 0
	//	}
	//
	// aExt, bExt := filepath.Ext(a.Path), filepath.Ext(b.Path)
	//
	//	if aExt == bExt {
	//		if a.Path > b.Path {
	//			return 1
	//		} else {
	//			return -1
	//		}
	//	}
	//
	//	if aExt > bExt {
	//		return 1
	//	} else {
	//
	//		return -1
	//	}
}

func (r *TorrentEntity) FlatFiles() []*FileEntity {
	return slices.SortedFunc(maps.Values(r.FilePathMap), fileCompare)
}

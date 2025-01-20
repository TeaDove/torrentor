package schemas

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/conv_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"torrentor/services/ffmpeg_service"
)

type FileEntity struct {
	Name string `json:"name"`
	Path string `json:"path"`

	Mimetype string          `json:"mimetype,omitempty"`
	Size     conv_utils.Byte `json:"size"`

	Completed bool `json:"completed"`

	Meta    ffmpeg_service.Metadata `json:"meta,omitempty"`
	Obj     *torrent.File           `json:"-"`
	Torrent *TorrentEntity          `json:"-"`
}

func (r *FileEntity) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("hash", redact_utils.Trim(r.Hash())).
		Str("name", r.Name).
		Str("size", r.Size.String())
}

func (r *FileEntity) Hash() string {
	if r.Obj != nil && r.Obj.FileInfo().Sha1 != "" {
		return r.Obj.FileInfo().Sha1
	}

	hasher := sha1.New()
	hasher.Write([]byte(r.Path))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
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
	return path.Join(r.Torrent.Location(), r.Path)
}

func (r *FileEntity) LocationInUnpack() string {
	return path.Join(r.Torrent.LocationInUnpack(), r.Hash())
}

func (r *FileEntity) LocationInUnpackAsStream(stream *ffmpeg_service.Stream, suffix string) string {
	return path.Join(r.Torrent.LocationInUnpack(), r.Hash(), stream.StreamFile(suffix))
}

type FileWithContent struct {
	*FileEntity
	OSFile *os.File `json:"-"`
}

func fileCompare(a *FileEntity, b *FileEntity) int {
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

func (r *TorrentEntity) FlatFiles() []*FileEntity {
	return slices.SortedFunc(maps.Values(r.FilePathMap), fileCompare)
}

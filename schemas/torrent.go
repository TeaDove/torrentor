package schemas

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/teadove/teasutils/utils/conv_utils"
	"path"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
)

type TorrentEntity struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`

	FilePathMap map[string]*FileEntity `json:"filePathMap,omitempty"`
	FileHashMap map[string]*FileEntity `json:"-"`

	Size conv_utils.Byte `json:"size"`

	InfoHash  metainfo.Hash `json:"infoHash"`
	Completed bool          `json:"completed"`

	Meta Meta `json:"meta,omitempty"`

	TorrentDataDir string           `json:"-"`
	UnpackDataDir  string           `json:"-"`
	Obj            *torrent.Torrent `json:"-"`
}

type Meta struct {
	Pieces      uint64 `json:"pieces"`
	PieceLength uint64 `json:"piecesLength"`
	Magnet      string `json:"magnet"`
}

func (r *TorrentEntity) AppendFile(file *FileEntity) {
	r.Size += file.Size

	r.FilePathMap[file.Path] = file
	r.FileHashMap[file.Hash()] = file
}

func (r *TorrentEntity) Location() string {
	return path.Join(r.TorrentDataDir, r.InfoHash.String(), r.Name)
}

func (r *TorrentEntity) LocationInUnpack() string {
	return path.Join(r.UnpackDataDir, r.InfoHash.String(), r.Name)
}

func (r *TorrentEntity) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("infohash", redact_utils.Trim(r.InfoHash.String())).
		Str("name", r.Name).
		Str("size", r.Size.String())
}

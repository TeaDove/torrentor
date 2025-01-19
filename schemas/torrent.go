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

	Files map[string]*FileEntity `json:"files,omitempty"`

	Size conv_utils.Byte `json:"size"`

	InfoHash  metainfo.Hash `json:"infoHash"`
	Completed bool          `json:"completed"`

	Meta Meta `json:"meta,omitempty"`
}

type Meta struct {
	Pieces      uint64 `json:"pieces"`
	PieceLength uint64 `json:"piecesLength"`
	Magnet      string `json:"magnet"`
}

func (r *TorrentEntity) AppendFile(file *FileEntity) {
	r.Size += file.Size

	r.Files[file.Path] = file
}

func (r *TorrentEntity) Location(dataDir string) string {
	return path.Join(dataDir, r.InfoHash.String())
}

func (r *TorrentEntity) FileLocation(dataDir string, filePath string) string {
	return path.Join(r.Location(dataDir), r.Name, filePath)
}

func (r *TorrentEntity) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("infohash", redact_utils.Trim(r.InfoHash.String())).
		Str("name", r.Name).
		Str("size", r.Size.String())
}

type TorrentEntityPop struct {
	*TorrentEntity
	Obj *torrent.Torrent `json:"-"`
}

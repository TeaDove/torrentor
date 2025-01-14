package schemas

import (
	"github.com/teadove/teasutils/utils/conv_utils"
	"path"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
)

type TorrentEntity struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name" gorm:"nindex"`

	// TODO remove
	Files map[string]FileEntity `json:"files,omitempty" gorm:"-:all"`

	Size conv_utils.Byte `json:"size"`

	InfoHash  string `json:"infoHash" gorm:"unique"`
	Completed bool   `json:"completed"`

	Meta Meta `json:"meta,omitempty" gorm:"embedded;embeddedPrefix:meta_"`
}

type Meta struct {
	Pieces      uint64 `json:"pieces,omitempty"`
	PieceLength uint64 `json:"piecesLength,omitempty"`
	Magnet      string `json:"magnet"`
}

func (r *TorrentEntity) AppendFile(file FileEntity) {
	r.Size += file.Size

	r.Files[file.Path] = file
}

func (r *TorrentEntity) Location(dataDir string) string {
	return path.Join(dataDir, r.InfoHash)
}

func (r *TorrentEntity) FileLocation(dataDir string, filePath string) string {
	return path.Join(r.Location(dataDir), r.Name, filePath)
}

func (r *TorrentEntity) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("id", r.ID.String()).
		Str("name", r.Name).
		Str("magnet", redact_utils.Trim(r.Meta.Magnet)).
		Str("size", r.Size.String())
}

type TorrentEntityPop struct {
	TorrentEntity
	Obj *torrent.Torrent `json:"-"`
}

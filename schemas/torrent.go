package schemas

import (
	"github.com/anacrolix/torrent"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/redact_utils"
	"path"
	"time"
)

type TorrentEntity struct {
	ID      uuid.UUID `json:"id"`
	AddedAt time.Time `json:"createdAt"`
	Name    string    `json:"name"`

	Root  FileEntity            `json:"root,omitempty"`
	Files map[string]FileEntity `json:"files,omitempty"`

	InfoHash  string `json:"infoHash"`
	Completed bool   `json:"completed"`

	Meta Meta `json:"meta,omitempty"`
}

type Meta struct {
	Pieces      uint64 `json:"pieces,omitempty"`
	PieceLength uint64 `json:"piecesLength,omitempty"`
	Magnet      string `json:"magnet"`
}

func (r *TorrentEntity) AppendFile(file FileEntity) {
	r.Files[file.Path] = file
}

func (r *TorrentEntity) Location(dataDir string) string {
	return path.Join(dataDir, r.InfoHash)
}

func (r *TorrentEntity) FileLocation(dataDir string, filePath string) string {
	return path.Join(r.Location(dataDir), filePath)
}

func (r *TorrentEntity) ZerologDict() *zerolog.Event {
	return zerolog.Dict().
		Str("id", r.ID.String()).
		Str("name", r.Name).
		Str("magnet", redact_utils.Trim(r.Meta.Magnet)).
		Str("size", r.Root.SizeRepr)
}

type TorrentEntityPop struct {
	TorrentEntity
	Obj *torrent.Torrent `json:"-"`
}

package settings

import (
	"time"

	"github.com/teadove/teasutils/utils/settings_utils"
)

type baseSettings struct {
	StartedAt time.Time

	UnpackDataDir string `env:"UNPACK_DATA_DIR" envDefault:"./data/unpack/"`
	DataDir       string `env:"DATA_DIR"        envDefault:"./data/torrent/"`

	URL string `env:"URL" envDefault:"0.0.0.0:8080"`
}

// Settings
// nolint: gochecknoglobals // need it
var Settings = settings_utils.MustGetSetting[baseSettings]("TORRENTOR_")

func init() { //nolint:gochecknoinits // required for started at
	Settings.StartedAt = time.Now().UTC()
}

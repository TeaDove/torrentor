package settings

import (
	"time"

	"github.com/teadove/teasutils/utils/settings_utils"
)

type tgSettings struct {
	BotToken string `env:"BOT_TOKEN" envDefault:"BAD_TOKEN"`
}

type webServerSettings struct {
	URL         string `env:"URL"          envDefault:"0.0.0.0:8081"`
	ExternalURL string `env:"EXTERNAL_URL" envDefault:"http://127.0.0.1:8081"`
}

type torrentSettings struct {
	DataDir string `env:"DATA_DIR" envDefault:"./data/torrent/"`
}

type baseSettings struct {
	UnpackDataDir string `env:"UNPACK_DATA_DIR" envDefault:"./data/unpack/"`
	StartedAt     time.Time

	TG        tgSettings        `envPrefix:"TG__"`
	WebServer webServerSettings `envPrefix:"WEB__"`
	Torrent   torrentSettings   `envPrefix:"TORRENT__"`
}

// Settings
// nolint: gochecknoglobals // need it
var Settings = settings_utils.MustGetSetting[baseSettings]("TORRENTOR_")

func init() { //nolint:gochecknoinits // required for started at
	Settings.StartedAt = time.Now()
}

package settings

import (
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/settings_utils"
)

type tgSettings struct {
	BotToken string `env:"BOT_TOKEN" envDefault:"BAD_TOKEN"`
}

type baseSettings struct {
	TG tgSettings `envPrefix:"TG__"`
}

var Settings = settings_utils.MustInitSetting[baseSettings](logger_utils.NewLoggedCtx(), "TORRENTO_", "TG.BotToken")

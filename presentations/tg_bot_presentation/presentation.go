package tg_bot_presentation

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/must_utils"
	"strings"
	"sync"
	"torrentor/services/torrentor_service"
)

type Presentation struct {
	bot *tgbotapi.BotAPI

	torrentorService *torrentor_service.Service
}

func NewPresentation(
	_ context.Context,
	bot *tgbotapi.BotAPI,
	torrentorService *torrentor_service.Service,
) (*Presentation, error) {
	command := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "download",
			Description: "Скачать",
		},
		tgbotapi.BotCommand{
			Command:     "stats",
			Description: "Статистика",
		},
	)
	_, err := bot.Request(command)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set commands")
	}

	return &Presentation{bot: bot, torrentorService: torrentorService}, nil
}

func (r *Presentation) PollerRun(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	// TODO move to settings
	u.Timeout = 10
	updates := r.bot.GetUpdatesChan(u)

	zerolog.Ctx(ctx).Info().Msg("bot.polling.started")

	var wg sync.WaitGroup

	for update := range updates {
		wg.Add(1)

		go must_utils.DoOrLogWithStacktrace(
			func(ctx context.Context) error {
				defer func() {
					err := must_utils.AnyToErr(recover())
					if err == nil {
						return
					}

					zerolog.Ctx(ctx).
						Err(err).Stack().
						Interface("update", update).
						Msg("panic.in.process.update")
				}()
				return r.processUpdate(ctx, &wg, &update)
			},
			"error.during.update.process",
		)(ctx)
	}

	wg.Wait()
}

func extractCommandAndText(text string, botUsername string, isChat bool) (string, string) {
	// TODO move to other module
	if len(text) <= 1 || text[0] != '/' || strings.HasPrefix(text, "/@") {
		return "", text
	}

	spaceIdx := strings.Index(text, " ")
	atIdx := strings.Index(text, "@")
	if atIdx == -1 && isChat {
		return "", text
	}

	if atIdx != -1 && (spaceIdx == -1 || spaceIdx > atIdx) {
		var extractedUsername string
		if spaceIdx == -1 {
			extractedUsername = text[atIdx:]
		} else {
			extractedUsername = text[atIdx:spaceIdx]
		}

		if extractedUsername == fmt.Sprintf("@%s", botUsername) {
			if spaceIdx == -1 {
				return text[1:atIdx], ""
			}
			return text[1:atIdx], text[spaceIdx+1:]
		} else {
			return "", text
		}
	}

	if spaceIdx == -1 {
		return text[1:], ""
	} else {
		return text[1:spaceIdx], text[spaceIdx+1:]
	}

}

func (r *Presentation) processUpdate(ctx context.Context, wg *sync.WaitGroup, update *tgbotapi.Update) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer wg.Done()
	c := r.makeCtx(ctx, update)

	zerolog.Ctx(c.ctx).Debug().Msg("processing.update")

	// TODO set advected commands
	switch c.command {
	case "download":
		c.tryReplyOnErr(c.Download())
	case "stats":
		c.tryReplyOnErr(c.Stats())
	}

	return nil
}

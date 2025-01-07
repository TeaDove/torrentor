package tg_bot_presentation

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func extractCommand(text string) string {
	if len(text) < 2 || text[0] != '/' {
		return ""
	}
	idx := strings.Index(text, " ")
	var command string
	if idx == -1 {
		command = text[1:]
	} else {
		command = text[1:idx]
	}

	return command
}

func (r *Presentation) processUpdate(ctx context.Context, wg *sync.WaitGroup, update *tgbotapi.Update) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer wg.Done()
	c := r.makeCtx(ctx, update)

	if c.command == "" {
		zerolog.Ctx(c.ctx).Debug().Msg("processing.update")
	}

	switch c.command {
	case "download":
		c.tryReplyOnErr(c.Download())
	case "stats":
		c.tryReplyOnErr(c.Stats())
	}

	return nil
}

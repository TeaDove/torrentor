package tg_bot_presentation

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/must_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
	"strconv"
	"sync"
	"time"
	"torrentor/suppliers/torrent_supplier"
)

type Presentation struct {
	bot *tgbotapi.BotAPI

	torrentSupplier *torrent_supplier.Supplier
}

func NewPresentation(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	torrentSupplier *torrent_supplier.Supplier,
) (*Presentation, error) {
	return &Presentation{bot: bot, torrentSupplier: torrentSupplier}, nil
}

func (r *Presentation) PollerRun(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	// TODO move to settings
	u.Timeout = 10
	updates := r.bot.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		wg.Add(1)
		go must_utils.DoOrLogWithStacktrace(
			func(ctx context.Context) error {
				innerCtx, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()

				return r.processUpdate(innerCtx, &wg, &update)
			},
			"error.during.update.process",
		)(ctx)
	}

	wg.Wait()
}

func (r *Presentation) processUpdate(ctx context.Context, wg *sync.WaitGroup, update *tgbotapi.Update) error {
	defer wg.Done()
	chat := update.FromChat()
	if chat != nil && chat.Title != "" {
		ctx = logger_utils.WithStrContextLog(ctx, "chat_title", chat.Title)
	}

	if update.Message != nil {
		ctx = logger_utils.WithStrContextLog(ctx, "message_id", strconv.Itoa(update.Message.MessageID))
		if update.Message.Text != "" {
			ctx = logger_utils.WithStrContextLog(ctx, "message_text", redact_utils.Trim(update.Message.Text))
		}
	}

	user := update.SentFrom()
	if user != nil {
		ctx = logger_utils.WithStrContextLog(ctx, "user", user.String())
	}

	zerolog.Ctx(ctx).Debug().Msg("processing.update")
	err := r.Download(ctx, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *Presentation) reply(update *tgbotapi.Update, format string, a ...any) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(format, a...))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := r.bot.Send(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

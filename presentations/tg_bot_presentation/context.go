package tg_bot_presentation

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/redact_utils"
	"strconv"
)

type Context struct {
	presentation *Presentation
	ctx          context.Context

	update   tgbotapi.Update
	text     string
	fulltext string
	command  string

	user *tgbotapi.User
	chat *tgbotapi.Chat
}

//func (r *Context) WithContext(ctx context.Context) *Context {
//	r.ctx = ctx
//	return r
//}

func (r *Context) addLogCtx() {
	if r.chat != nil && r.chat.Title != "" {
		r.ctx = logger_utils.WithStrContextLog(r.ctx, "chat_title", r.chat.Title)
	}

	if r.update.Message != nil {
		r.ctx = logger_utils.WithStrContextLog(r.ctx, "message_id", strconv.Itoa(r.update.Message.MessageID))
		if r.text != "" {
			r.ctx = logger_utils.WithStrContextLog(r.ctx, "message_text", redact_utils.Trim(r.text))
		}
	}

	if r.user != nil {
		r.ctx = logger_utils.WithStrContextLog(r.ctx, "user", r.user.String())
	}
	if r.command != "" {
		r.ctx = logger_utils.WithStrContextLog(r.ctx, "command", r.command)
	}
}

func (r *Presentation) makeCtx(ctx context.Context, update *tgbotapi.Update) Context {
	c := Context{
		presentation: r,
		update:       *update,
		chat:         update.FromChat(),
		user:         update.SentFrom(),
	}

	if update.Message != nil {
		c.fulltext = update.Message.Text
	}

	c.ctx = ctx
	c.command = extractCommand(c.fulltext)
	if c.command == "" {
		c.text = c.fulltext
	} else if len(c.command)+2 < len(c.fulltext) {
		c.text = c.fulltext[len(c.command)+2:]
	}

	c.addLogCtx()

	return c
}

func (r *Context) Log() *zerolog.Logger {
	return zerolog.Ctx(r.ctx)
}

func (r *Context) LogWithUpdate() zerolog.Logger {
	return zerolog.Ctx(r.ctx).With().Interface("update", r.update).Logger()
}

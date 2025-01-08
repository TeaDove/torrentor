package tg_bot_presentation

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Context) reply(format string, a ...any) error {
	_, err := r.replyWithMessage(format, a...)
	return err
}

func (r *Context) tryReply(format string, a ...any) {
	_, err := r.replyWithMessage(format, a...)
	if err != nil {
		r.Log().
			Error().Stack().
			Err(err).
			Str("format", format).
			Interface("args", a).
			Msg("failed.to.reply")
	}

}

func (r *Context) replyWithMessage(format string, a ...any) (tgbotapi.Message, error) {
	msgReq := tgbotapi.NewMessage(r.update.Message.Chat.ID, fmt.Sprintf(format, a...))
	msgReq.ReplyToMessageID = r.update.Message.MessageID
	msgReq.ParseMode = tgbotapi.ModeHTML

	msg, err := r.presentation.bot.Send(msgReq)
	if err != nil {
		return tgbotapi.Message{}, errors.Wrap(err, "failed to send message")
	}

	return msg, nil
}

func (r *Context) editMsgText(msg *tgbotapi.Message, format string, a ...any) error {
	editMessageTextReq := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, fmt.Sprintf(format, a...))
	editMessageTextReq.ParseMode = tgbotapi.ModeHTML

	_, err := r.presentation.bot.Send(editMessageTextReq)
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	return nil
}

func (r *Context) tryReplyOnErr(err error) {
	if err == nil {
		return
	}

	zerolog.Ctx(r.ctx).Error().Stack().Err(err).Msg("unexpected.error")
	err = r.reply("Unexpected error occurred: %s", err.Error())
	if err != nil {
		zerolog.Ctx(r.ctx).Error().Stack().Err(err).Msg("failed.to.try.reply.on.err")
	}
}

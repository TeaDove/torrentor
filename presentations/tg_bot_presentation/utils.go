package tg_bot_presentation

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (r *Context) reply(text string) error {
	_, err := r.replyWithMessage(text)
	return err
}

func (r *Context) tryReply(text string) {
	_, err := r.replyWithMessage(text)
	if err != nil {
		r.Log().
			Error().Stack().
			Err(err).
			Str("text", text).
			Msg("failed.to.reply")
	}
}

func (r *Context) replyWithMessage(text string) (tgbotapi.Message, error) {
	msgReq := tgbotapi.NewMessage(r.update.Message.Chat.ID, text)
	msgReq.ReplyToMessageID = r.update.Message.MessageID
	msgReq.ParseMode = tgbotapi.ModeHTML

	msg, err := r.presentation.bot.Send(msgReq)
	if err != nil {
		return tgbotapi.Message{}, errors.Wrap(err, "failed to send message")
	}

	return msg, nil
}

func (r *Context) editMsgText(msg *tgbotapi.Message, text string) error {
	editMessageTextReq := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, text)
	editMessageTextReq.ParseMode = tgbotapi.ModeHTML

	newMsg, err := r.presentation.bot.Send(editMessageTextReq)
	if err != nil {
		return errors.Wrap(err, "failed to edit message")
	}

	*msg = newMsg

	return nil
}

func (r *Context) tryReplyOnErr(err error) {
	if err == nil {
		return
	}

	zerolog.Ctx(r.ctx).Error().Stack().Err(err).Msg("unexpected.error")

	err = r.reply(fmt.Sprintf("Unexpected error occurred: %s", err.Error()))
	if err != nil {
		zerolog.Ctx(r.ctx).Error().Stack().Err(err).Msg("failed.to.try.reply.on.err")
	}
}

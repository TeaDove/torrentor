package tg_bot_presentation

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func (r *Presentation) Download(ctx context.Context, update *tgbotapi.Update) error {
	link := update.Message.Text

	msg, err := r.replyWithMessage(update, "Пытаемся получить метадату торрента")
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	torrent, statsChan, err := r.torrentorService.DownloadAndSaveFromMagnet(ctx, link)
	if err != nil {
		return errors.Wrap(err, "failed to download magnet")
	}

	r.bot.Send
	err = r.reply(update, "%s", torrent.Name)
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	err = r.reply(update, "http://localhost:8081/torrent/%s", torrent.Id)
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	for stats := range statsChan {
		err = r.reply(
			update,
			"Peers: %d/%d\nPieces complete: %d\nRead: %f MB\nWritten: %f MB",
			stats.ActivePeers,
			stats.TotalPeers,
			stats.PiecesComplete,
			converters_utils.ToMegaByte(stats.BytesRead.Int64()),
			converters_utils.ToMegaByte(stats.BytesWritten.Int64()),
		)
		if err != nil {
			return errors.Wrap(err, "failed to send reply")
		}
	}

	err = r.reply(update, "Done!")
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	return nil
}

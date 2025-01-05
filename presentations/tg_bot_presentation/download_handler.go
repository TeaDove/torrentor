package tg_bot_presentation

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (r *Presentation) Download(ctx context.Context, update *tgbotapi.Update) error {
	link := update.Message.Text

	torrentPath, err := r.torrentSupplier.DownloadMagnet(ctx, link)
	if err != nil {
		return errors.Wrap(err, "failed to download magnet")
	}

	err = r.reply(update, "%s", torrentPath)
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	fileID := ""

	err = r.reply(update, "%s", fmt.Sprintf("http://localhost:8080/%s", fileID))
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	http.Handle(fmt.Sprintf("/%s", fileID), http.FileServer(http.Dir(path.Join("./data/torrent", torrentPath))))

	//nolint: gosec // TODO move to settings
	if err = http.ListenAndServe(":8080", nil); err != nil {
		// fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	err = r.reply(update, "%s", torrentPath)
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	return nil
}

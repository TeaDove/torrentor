package tg_bot_presentation

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
	"time"
)

const torrentDownloadingTmpl = `<code>%s</code>
Доступно <a href="%s">тут</a>

Статус:
%%s
`

func (r *Context) Download() error {
	link := r.text

	msg, err := r.replyWithMessage("Пытаемся получить метадату торрента...")
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	t0 := time.Now()
	torrent, statsChan, err := r.presentation.torrentorService.DownloadAndSaveFromMagnet(r.ctx, link)
	if err != nil {
		return errors.Wrap(err, "failed to download magnet")
	}

	msgTextTmpl := fmt.Sprintf(torrentDownloadingTmpl, torrent.Name, "https://example.com/torrent/"+torrent.Id.String())
	err = r.editMsgText(&msg, msgTextTmpl, "Подключаемся")
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	for stats := range statsChan {
		err = r.editMsgText(
			&msg,
			msgTextTmpl,
			fmt.Sprintf(
				"Peers: %d/%d\nPieces complete: %d/%d\nRead: %f MB\nWritten: %f MB\nElapsed time: %s",
				stats.ActivePeers,
				stats.TotalPeers,
				stats.PiecesComplete,
				torrent.Pieces,
				converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesRead.Int64()), 1),
				converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesWritten.Int64()), 1),
				time.Since(t0).String(),
			),
		)
		if err != nil {
			return errors.Wrap(err, "failed to send reply")
		}
	}

	err = r.editMsgText(&msg, msgTextTmpl, "Done!")
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	return nil
}

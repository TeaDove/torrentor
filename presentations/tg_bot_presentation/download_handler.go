package tg_bot_presentation

import (
	"fmt"
	"strings"
	"time"
	"torrentor/settings"

	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
)

const torrentDownloadingTmpl = `<code>%s</code>
Доступно <a href="%s">тут</a>

Статус:
%%s
`

func (r *Context) Download() error {
	link := r.text
	if !strings.HasPrefix(link, "magnet:?xt=urn") {
		err := r.reply("Magnet link required as argument, for example")
		if err != nil {
			return err
		}

		err = r.reply(
			"/download magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10&dn=Sintel&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&tr=wss%3A%2F%2Ftracker.fastcast.nz&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F",
		)
		if err != nil {
			return err
		}

		return nil
	}

	msg, err := r.replyWithMessage("Пытаемся получить метадату торрента...")
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	t0 := time.Now()

	torrent, statsChan, err := r.presentation.torrentorService.DownloadAndSaveFromMagnet(r.ctx, link)
	if err != nil {
		return errors.Wrap(err, "failed to download magnet")
	}

	msgTextTmpl := fmt.Sprintf(
		torrentDownloadingTmpl,
		torrent.Name,
		settings.Settings.WebServer.ExternalURL+"/torrents/"+torrent.ID.String(),
	)

	err = r.editMsgText(&msg, fmt.Sprintf(msgTextTmpl, "Подключаемся"))
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	for stats := range statsChan {
		err = r.editMsgText(
			&msg,
			fmt.Sprintf(
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
			),
		)
		if err != nil {
			return errors.Wrap(err, "failed to send reply")
		}
	}

	err = r.editMsgText(&msg, fmt.Sprintf(msgTextTmpl, "Готово!"))
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	return nil
}

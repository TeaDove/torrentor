package tg_bot_presentation

import (
	"fmt"
	"strings"
	"time"
	"torrentor/settings"

	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
)

const torrentDownloadingTmpl = `<a href="%s">%s</a>

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

	msg, err := r.replyWithMessage("Getting torrent metadata")
	if err != nil {
		return errors.Wrap(err, "failed to reply")
	}

	torrent, statsChan, err := r.presentation.torrentorService.DownloadAndSaveFromMagnet(r.ctx, link)
	if err != nil {
		return errors.Wrap(err, "failed to download magnet")
	}

	msgTextTmpl := fmt.Sprintf(
		torrentDownloadingTmpl,
		settings.Settings.WebServer.ExternalURL+"/torrents/"+torrent.ID.String(),
		torrent.Name,
	)

	err = r.editMsgText(&msg, fmt.Sprintf(msgTextTmpl, "Подключаемся"))
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	t0 := time.Now()

	for stats := range statsChan {
		bytesDone := uint64(stats.PiecesComplete) * torrent.Meta.PieceLength

		err = r.editMsgText(
			&msg,
			fmt.Sprintf(
				msgTextTmpl,
				fmt.Sprintf(
					"Peers: %d / %d\nComplete: %s / %s\nSpeed: %s/s",
					stats.ActivePeers,
					stats.TotalPeers,
					converters_utils.ToClosestByteAsString(bytesDone, 2),
					converters_utils.ToClosestByteAsString(torrent.Meta.Pieces*torrent.Meta.PieceLength, 2),
					converters_utils.ToClosestByteAsString(float64(bytesDone)/(time.Since(t0).Seconds()), 2),
				),
			),
		)
		if err != nil {
			return errors.Wrap(err, "failed to send reply")
		}
	}

	err = r.editMsgText(&msg, fmt.Sprintf(msgTextTmpl, "Done!"))
	if err != nil {
		return errors.Wrap(err, "failed to send reply")
	}

	return nil
}

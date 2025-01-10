package tg_bot_presentation

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func makeStatsMsgText(stats *torrent.ClientStats) string {
	return fmt.Sprintf(
		"Peers: %d\nRead: %f MB\nWritten: %f MB",
		stats.NumPeersDialedSuccessfullyAfterHolepunchConnect,
		converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesRead.Int64()), 1),
		converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesWritten.Int64()), 1),
	)
}

func (r *Context) Stats() error {
	statsChan := r.presentation.torrentorService.Stats(r.ctx)
	stats, ok := <-statsChan
	if !ok {
		r.tryReply("Stats channel closed, this can occur because torrent client is not ready yet")
	}

	msg, err := r.replyWithMessage(makeStatsMsgText(&stats))
	if err != nil {
		return errors.Wrap(err, "failed to reply with msg")
	}

	for stats = range statsChan {
		text := makeStatsMsgText(&stats)
		if text == msg.Text {
			continue
		}

		err = r.editMsgText(&msg, text)
		if err != nil {
			return errors.Wrap(err, "failed to reply with msg")
		}
	}

	err = r.editMsgText(&msg, msg.Text+"\n\nStats streaming stopped, for more - send /stats")
	if err != nil {
		return errors.Wrap(err, "failed to reply with msg")
	}

	return nil
}

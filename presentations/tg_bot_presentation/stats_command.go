package tg_bot_presentation

import (
	"fmt"
	"time"

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
	serviceStats, statsChan, err := r.presentation.torrentorService.Stats(r.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get stats")
	}

	err = r.reply(fmt.Sprintf(`
Since: %s (%s)
Torrents count: %d (files: %d)
Total size: %s
`,
		serviceStats.StartedAt,
		time.Since(serviceStats.StartedAt),
		serviceStats.TorrentsCount,
		serviceStats.FilesCount,
		converters_utils.ToClosestByteAsString(serviceStats.TotalSize, 2),
	),
	)
	if err != nil {
		return errors.Wrap(err, "failed to send stats message")
	}

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

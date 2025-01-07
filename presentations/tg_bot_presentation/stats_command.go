package tg_bot_presentation

import (
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func (r *Context) Stats() error {
	statsChan := r.presentation.torrentorService.Stats()
	stats, ok := <-statsChan
	if !ok {
		return errors.New("stats channel closed")
	}

	msg, err := r.replyWithMessage("Peers: %d\nRead: %f MB\nWritten: %f MB",
		stats.NumPeersDialedSuccessfullyAfterHolepunchConnect,
		converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesRead.Int64()), 1),
		converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesWritten.Int64()), 1),
	)
	if err != nil {
		return errors.Wrap(err, "failed to reply with msg")
	}

	for stats = range statsChan {
		err = r.editMsgText(&msg, "Peers: %d\nRead: %f MB\nWritten: %f MB",
			stats.NumPeersDialedSuccessfullyAfterHolepunchConnect,
			converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesRead.Int64()), 1),
			converters_utils.ToFixed(converters_utils.ToMegaByte(stats.BytesWritten.Int64()), 1),
		)
		if err != nil {
			return errors.Wrap(err, "failed to reply with msg")
		}
	}

	return nil
}

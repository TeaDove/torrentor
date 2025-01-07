package torrent_supplier

import (
	"context"
	"github.com/anacrolix/torrent"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/converters_utils"
	"time"

	"github.com/pkg/errors"
)

func (r *Supplier) ExportStats(ctx context.Context, t *torrent.Torrent) <-chan torrent.TorrentStats {
	torrentStats := make(chan torrent.TorrentStats, 1)
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				torrentStats <- t.Stats()

			case <-t.Complete().On():
				close(torrentStats)
				return

			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return torrentStats
}

func (r *Supplier) AddMagnetAndGetInfoAndStartDownload(
	ctx context.Context,
	magnetLink string,
) (*torrent.Torrent, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t, err := r.client.AddMagnet(magnetLink)
	if err != nil {
		return t, errors.Wrap(err, "failed to add magnet")
	}

	<-t.GotInfo()
	t.DownloadAll()

	zerolog.Ctx(ctx).
		Info().
		Str("torrent_name", t.Name()).
		Float64("total_size_mb", converters_utils.ToMegaByte(t.Info().TotalLength())).
		Msg("torrent.info.got")

	return t, nil
}

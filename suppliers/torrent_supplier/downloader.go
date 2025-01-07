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
	torrentStats := make(chan torrent.TorrentStats)
	ticker := time.NewTicker(time.Second)

	go func() {
		defer close(torrentStats)
		// possible mem lick
		for {
			select {
			case <-ticker.C:
				torrentStats <- t.Stats()

			case <-t.Complete().On():
				zerolog.Ctx(ctx).Info().Str("torrent", t.Name()).Msg("download.complete")
				return
			case <-ctx.Done():
				zerolog.Ctx(ctx).Info().Str("torrent", t.Name()).Msg("ctx.canceled.but.not.download")
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

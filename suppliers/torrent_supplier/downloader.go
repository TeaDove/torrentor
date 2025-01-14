package torrent_supplier

import (
	"context"
	"github.com/teadove/teasutils/utils/conv_utils"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/rs/zerolog"

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
		Str("name", t.Name()).
		Str("size", conv_utils.ClosestByte(t.Info().TotalLength())).
		Msg("torrent.info.ready")

	return t, nil
}

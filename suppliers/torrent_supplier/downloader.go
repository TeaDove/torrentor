package torrent_supplier

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
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
	torrentObj, err := r.client.AddMagnet(magnetLink)
	if err != nil {
		return torrentObj, errors.Wrap(err, "failed to add magnet")
	}

	<-torrentObj.GotInfo()
	torrentObj.DownloadAll()

	zerolog.Ctx(ctx).
		Info().
		Str("name", torrentObj.Name()).
		Str("size", conv_utils.ClosestByte(torrentObj.Info().TotalLength())).
		Msg("torrent.info.ready")

	return torrentObj, nil
}

func (r *Supplier) GetTorrentByInfoHash(
	ctx context.Context,
	infoHash metainfo.Hash,
) (*torrent.Torrent, error) {
	torrentObj, _ := r.client.AddTorrentInfoHash(infoHash)
	select {
	case <-torrentObj.GotInfo():
		return torrentObj, nil
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "failed to get torrent info")
	}
}

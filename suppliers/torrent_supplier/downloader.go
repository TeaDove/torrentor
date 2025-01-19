package torrent_supplier

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
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

func waitForInfo(ctx context.Context, torrentObj *torrent.Torrent) error {
	select {
	case <-torrentObj.GotInfo():
		return nil
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "failed to get torrent info")
	}
}

func (r *Supplier) AddMagnetAndGetInfoAndStartDownload(
	ctx context.Context,
	magnetLink string,
) (*torrent.Torrent, error) {
	torrentObj, err := r.client.AddMagnet(magnetLink)
	if err != nil {
		return torrentObj, errors.Wrap(err, "failed to add magnet")
	}

	err = waitForInfo(ctx, torrentObj)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait for magnet")
	}

	torrentObj.DownloadAll()

	return torrentObj, nil
}

func (r *Supplier) GetTorrentByInfoHash(
	ctx context.Context,
	infoHash metainfo.Hash,
) (*torrent.Torrent, error) {
	torrentObj, _ := r.client.AddTorrentInfoHash(infoHash)
	return torrentObj, waitForInfo(ctx, torrentObj)
}

package torrent_supplier

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/converters_utils"
)

func (r *Supplier) DownloadMagnet(ctx context.Context, magnetLink string) (string, error) {
	t, err := r.client.AddMagnet(magnetLink)
	if err != nil {
		return "", errors.Wrap(err, "failed to add magnet")
	}

	<-t.GotInfo()
	zerolog.Ctx(ctx).
		Info().
		Str("torrent_name", t.Name()).
		Float64("total_size_mb", converters_utils.ToMegaByte(t.Info().TotalLength())).
		Msg("torrent.info.got")

	t.DownloadAll()
	<-t.Complete().On()

	zerolog.Ctx(ctx).
		Info().
		Str("torrent_name", t.Name()).
		Float64("total_size_mb", converters_utils.ToMegaByte(t.Info().TotalLength())).
		Msg("torrent.downloaded")

	return t.Name(), nil
}

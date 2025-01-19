package torrentor_service

import (
	"context"
	"github.com/anacrolix/torrent/metainfo"
	"sync"
	"time"
	"torrentor/schemas"
	"torrentor/services/ffmpeg_service"
	"torrentor/suppliers/torrent_supplier"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/must_utils"
)

type Service struct {
	torrentSupplier *torrent_supplier.Supplier
	ffmpegService   *ffmpeg_service.Service

	hashToTorrent   map[metainfo.Hash]*schemas.TorrentEntityPop
	hashToTorrentMu sync.RWMutex

	torrentDataDir string
}

func NewService(
	ctx context.Context,
	torrentSupplier *torrent_supplier.Supplier,
	ffmpegService *ffmpeg_service.Service,
	scheduler *gocron.Scheduler,
	torrentDataDir string,
) (*Service, error) {
	r := &Service{
		torrentSupplier: torrentSupplier,
		ffmpegService:   ffmpegService,
		torrentDataDir:  torrentDataDir,
		hashToTorrent:   make(map[metainfo.Hash]*schemas.TorrentEntityPop, 10),
	}

	_, err := scheduler.
		//nolint: mnd // TODO move to settings
		Every(5*time.Minute).
		Do(must_utils.DoOrLog(r.restartDownloadForAllTorrents, "failed to restart download"), ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to schedule job")
	}

	return r, nil
}

package torrentor_service

import (
	"context"
	"os"
	"sync"
	"time"
	"torrentor/backend/schemas"
	"torrentor/backend/services/ffmpeg_service"
	"torrentor/backend/suppliers/torrent_supplier"

	"github.com/anacrolix/torrent/metainfo"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/must_utils"
)

type Service struct {
	torrentSupplier *torrent_supplier.Supplier
	ffmpegService   *ffmpeg_service.Service

	hashToTorrent   map[metainfo.Hash]*schemas.TorrentEntity
	hashToTorrentMu sync.RWMutex

	torrentDataDir string
	unpackDataDir  string
}

func NewService(
	ctx context.Context,
	torrentSupplier *torrent_supplier.Supplier,
	ffmpegService *ffmpeg_service.Service,
	scheduler *gocron.Scheduler,
	torrentDataDir string,
	unpackDataDir string,
) (*Service, error) {
	r := &Service{
		torrentSupplier: torrentSupplier,
		ffmpegService:   ffmpegService,
		torrentDataDir:  torrentDataDir,
		unpackDataDir:   unpackDataDir,
		hashToTorrent:   make(map[metainfo.Hash]*schemas.TorrentEntity, 10),
	}

	err := os.MkdirAll(r.unpackDataDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create unpack data directory")
	}

	_, err = scheduler.
		//nolint: mnd // TODO move to settings
		Every(5*time.Minute).
		Do(must_utils.DoOrLog(r.restartDownloadForAllTorrents, "failed to restart download"), ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to schedule job")
	}

	return r, nil
}

package torrentor_service

import (
	"context"
	"time"
	"torrentor/repositories/torrent_repository"
	"torrentor/services/ffmpeg_service"
	"torrentor/suppliers/torrent_supplier"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/teadove/teasutils/utils/must_utils"
)

type Service struct {
	torrentSupplier   *torrent_supplier.Supplier
	torrentRepository *torrent_repository.Repository
	ffmpegService     *ffmpeg_service.Service
}

func NewService(
	ctx context.Context,
	torrentSupplier *torrent_supplier.Supplier,
	torrentRepository *torrent_repository.Repository,
	ffmpegService *ffmpeg_service.Service,
	scheduler *gocron.Scheduler,
) (*Service, error) {
	r := &Service{
		torrentSupplier:   torrentSupplier,
		torrentRepository: torrentRepository,
		ffmpegService:     ffmpegService,
	}

	_, err := scheduler.
		//nolint: mnd // TODO move to settings
		Every(5*time.Minute).
		Do(must_utils.DoOrLog(r.restartDownload, "failed to restart download"), ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to schedule job")
	}

	return r, nil
}

package app_container

import (
	"context"
	"time"
	"torrentor/backend/presentations/api_presentation"
	"torrentor/backend/services/ffmpeg_service"
	"torrentor/backend/services/torrentor_service"
	"torrentor/backend/settings"
	"torrentor/backend/suppliers/torrent_supplier"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
)

type Container struct {
	WebPresentation *api_presentation.Presentation
}

func Build(ctx context.Context) (*Container, error) {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	torrentSupplier, err := torrent_supplier.NewSupplier(ctx, settings.Settings.DataDir)
	if err != nil {
		return nil, errors.Wrap(err, "could not create torrent supplier")
	}

	ffmpegService, err := ffmpeg_service.NewService(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ffmpeg service")
	}

	torrentorService, err := torrentor_service.NewService(
		ctx,
		torrentSupplier,
		ffmpegService,
		scheduler,
		settings.Settings.DataDir,
		settings.Settings.UnpackDataDir,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create torrent service")
	}

	webPresentation := api_presentation.NewPresentation(torrentorService)

	container := &Container{WebPresentation: webPresentation}

	return container, nil
}

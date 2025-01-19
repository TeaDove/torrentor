package app_container

import (
	"context"
	"time"
	"torrentor/presentations/tg_bot_presentation"
	"torrentor/presentations/web_app_presentation"
	"torrentor/repositories/torrent_repository"
	"torrentor/services/ffmpeg_service"
	"torrentor/services/torrentor_service"
	"torrentor/settings"
	"torrentor/suppliers/torrent_supplier"

	"github.com/go-co-op/gocron"
	"github.com/teadove/teasutils/utils/di_utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Container struct {
	TGBotPresentation *tg_bot_presentation.Presentation
	WebPresentation   *web_app_presentation.Presentation

	healths []di_utils.Health
	closers []di_utils.CloserWithContext
}

func (r *Container) Healths() []di_utils.Health {
	return r.healths
}

func (r *Container) Closers() []di_utils.CloserWithContext {
	return r.closers
}

func Build(ctx context.Context) (*Container, error) {
	scheduler := gocron.NewScheduler(time.UTC)

	torrentRepository, err := torrent_repository.NewRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create torrent repository")
	}

	// TODO move to settings
	torrentSupplier, err := torrent_supplier.NewSupplier(ctx, settings.Settings.Torrent.DataDir)
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
		torrentRepository,
		ffmpegService,
		scheduler,
		settings.Settings.Torrent.DataDir,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create torrent service")
	}

	bot, err := tgbotapi.NewBotAPI(settings.Settings.TG.BotToken)
	if err != nil {
		return nil, errors.Wrap(err, "could not create bot client")
	}

	tgBotPresentation, err := tg_bot_presentation.NewPresentation(ctx, bot, torrentorService)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tg_bot_presentation")
	}

	webPresentation, err := web_app_presentation.NewPresentation(ctx, torrentorService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create presentation")
	}

	scheduler.StartAsync()

	container := &Container{
		TGBotPresentation: tgBotPresentation,
		WebPresentation:   webPresentation,
		healths: []di_utils.Health{
			tgBotPresentation,
		},
		closers: []di_utils.CloserWithContext{
			torrentSupplier,
			tgBotPresentation,
			webPresentation,
		},
	}

	return container, nil
}

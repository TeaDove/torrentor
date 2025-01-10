package app_container

import (
	"context"
	"github.com/teadove/teasutils/utils/di_utils"
	"torrentor/infrastructure/buntdb_infrastructure"
	"torrentor/presentations/tg_bot_presentation"
	"torrentor/presentations/web_app_presentation"
	"torrentor/repositories/torrent_repository"
	"torrentor/services/torrentor_service"
	"torrentor/settings"
	"torrentor/suppliers/torrent_supplier"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Container struct {
	TGBotPresentation *tg_bot_presentation.Presentation
	WebPresentation   *web_app_presentation.Presentation

	healthCheckers []di_utils.Health
	stoppers       []di_utils.CloserWithContext
}

func (r *Container) Healths() []di_utils.Health {
	return r.healthCheckers
}

func (r *Container) Closers() []di_utils.CloserWithContext {
	return r.stoppers
}

func Build(ctx context.Context) (*Container, error) {
	db, err := buntdb_infrastructure.NewClientFromSettings()
	if err != nil {
		return nil, errors.Wrap(err, "opening db")
	}

	torrentRepository, err := torrent_repository.NewRepository(ctx, db, settings.Settings.Torrent.DataDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create torrent repository")
	}

	// TODO move to settings
	torrentSupplier, err := torrent_supplier.NewSupplier(ctx, settings.Settings.Torrent.DataDir)
	if err != nil {
		return nil, errors.Wrap(err, "could not create torrent supplier")
	}

	torrentorService, err := torrentor_service.NewService(ctx, torrentSupplier, torrentRepository)
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

	container := &Container{
		TGBotPresentation: tgBotPresentation,
		WebPresentation:   webPresentation,
		healthCheckers: []di_utils.Health{
			tgBotPresentation,
			torrentRepository,
		},
		stoppers: []di_utils.CloserWithContext{
			torrentSupplier,
			tgBotPresentation,
			torrentRepository,
			webPresentation,
		},
	}

	return container, nil
}

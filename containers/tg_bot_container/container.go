package tg_bot_container

import (
	"context"
	"torrentor/presentations/tg_bot_presentation"
	"torrentor/settings"
	"torrentor/suppliers/torrent_supplier"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Container struct {
	TGBotPresentation *tg_bot_presentation.Presentation

	healthCheckers []func(ctx context.Context) error
	stoppers       []func(ctx context.Context) error
}

func (r *Container) HealthCheckers() []func(ctx context.Context) error {
	return r.healthCheckers
}

func (r *Container) Stoppers() []func(ctx context.Context) error {
	return r.stoppers
}

func Build(ctx context.Context) (*Container, error) {
	// TODO move to settings
	torrentSupplier, err := torrent_supplier.NewSupplier(ctx, "./data/torrent/")
	if err != nil {
		return nil, errors.Wrap(err, "could not create torrent supplier")
	}

	bot, err := tgbotapi.NewBotAPI(settings.Settings.TG.BotToken)
	if err != nil {
		return nil, errors.Wrap(err, "could not create bot client")
	}

	tgBotPresentation, err := tg_bot_presentation.NewPresentation(ctx, bot, torrentSupplier)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tg_bot_presentation")
	}

	container := &Container{
		TGBotPresentation: tgBotPresentation,
		healthCheckers: []func(ctx context.Context) error{
			tgBotPresentation.Health,
		},
		stoppers: []func(ctx context.Context) error{
			torrentSupplier.Close,
			func(_ context.Context) error {
				bot.StopReceivingUpdates()
				return nil
			},
		},
	}

	return container, nil
}

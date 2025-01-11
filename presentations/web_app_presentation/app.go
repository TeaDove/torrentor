package web_app_presentation

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"net/http"
	"torrentor/presentations/web_app_presentation/views"
	"torrentor/services/torrentor_service"
	"torrentor/settings"
)

type Presentation struct {
	torrentorService *torrentor_service.Service
	fiberApp         *fiber.App
}

func NewPresentation(
	_ context.Context,
	torrentorService *torrentor_service.Service,
) (*Presentation, error) {
	app := fiber.New(fiber.Config{
		Immutable: true,
		Views:     html.NewFileSystem(http.FS(views.Static), ".html"),
	})
	app.Use(logCtxMiddleware())

	r := Presentation{
		torrentorService: torrentorService,
		fiberApp:         app,
	}

	r.fiberApp.Get("/", IndexForm)
	r.fiberApp.Get("/torrents/:id", r.TorrentForm)
	r.fiberApp.Get("/files/:id", r.FileForm)

	return &r, nil
}

func logCtxMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := logger_utils.AddLoggerToCtx(c.Context())
		ctx = logger_utils.WithStrContextLog(ctx,
			"app_method",
			fmt.Sprintf(
				"%s %s",
				c.Method(),
				c.Path(),
			),
		)
		ctx = logger_utils.WithStrContextLog(ctx, "ip", c.IP())
		c.SetContext(ctx)

		err := c.Next()
		if err != nil {
			zerolog.Ctx(ctx).Error().Stack().Err(err).Msg("error.in.request")
		}

		zerolog.Ctx(ctx).
			Debug().
			Int("status", c.Response().StatusCode()).
			Msg("request.completed")

		return err
	}
}

func (r *Presentation) Close(ctx context.Context) error {
	err := r.fiberApp.ShutdownWithContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error shutting down app")
	}

	return nil
}

func (r *Presentation) Run(_ context.Context) error {
	err := r.fiberApp.Listen(settings.Settings.WebServer.URL)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	return nil
}

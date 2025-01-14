package web_app_presentation

import (
	"context"
	"fmt"
	"net/http"
	"torrentor/presentations/web_app_presentation/views"
	"torrentor/schemas"
	"torrentor/services/torrentor_service"
	"torrentor/settings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
	"github.com/teadove/teasutils/utils/settings_utils"
)

type Presentation struct {
	torrentorService *torrentor_service.Service
	fiberApp         *fiber.App
}

func NewPresentation(
	_ context.Context,
	torrentorService *torrentor_service.Service,
) (*Presentation, error) {
	renderEngine := html.NewFileSystem(http.FS(views.Static), ".html")
	renderEngine.Funcmap["FileIsVideo"] = func(file schemas.FileEntity) bool {
		return file.IsVideo()
	}

	if !settings_utils.BaseSettings.Release {
		renderEngine.Debug(true)
		renderEngine.Reload(true)
	}

	app := fiber.New(fiber.Config{
		Immutable: true,
		Views:     renderEngine,
	})
	app.Use(logCtxMiddleware())

	r := Presentation{
		torrentorService: torrentorService,
		fiberApp:         app,
	}

	r.fiberApp.Get("/", IndexForm)
	r.fiberApp.Get("/torrents/:id", r.TorrentForm)
	r.fiberApp.Get("/torrents/:id/file", r.FileForm)
	r.fiberApp.Get("/torrents/:id/watch", r.WatchForm)

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
		c.SetContext(ctx)

		ctx = logger_utils.WithStrContextLog(ctx, "ip", c.IP())
		ctx = logger_utils.WithStrContextLog(ctx, "user_agent", c.Get("User-Agent"))

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

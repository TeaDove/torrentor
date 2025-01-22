package web_app_presentation

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"torrentor/presentations/web_app_presentation/views"
	"torrentor/schemas"
	"torrentor/services/torrentor_service"
	"torrentor/settings"

	"github.com/teadove/teasutils/utils/conv_utils"

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
	renderEngine.Funcmap["FileIsWatchable"] = func(file schemas.FileEntity) bool {
		return file.IsVideo() && file.Mimetype != schemas.MatroskaMimeType
	}
	renderEngine.Funcmap["FileIsVideo"] = func(file schemas.FileEntity) bool {
		return file.IsVideo()
	}
	renderEngine.Funcmap["FileAudioStreamsNames"] = func(file schemas.FileEntity) []string {
		return file.Meta.AudioStreamsAsStrings()
	}
	renderEngine.Funcmap["SizeRepr"] = func(size conv_utils.Byte) string {
		return size.String()
	}

	if !settings_utils.BaseSettings.Release {
		renderEngine.Reload(true)
	}

	app := fiber.New(fiber.Config{
		Immutable:    true,
		Views:        renderEngine,
		ErrorHandler: errHandler,
	})
	app.Use(logCtxMiddleware())

	r := Presentation{
		torrentorService: torrentorService,
		fiberApp:         app,
	}

	r.fiberApp.Get("/", IndexForm)
	r.fiberApp.Get("/torrents/:infohash", r.TorrentForm)
	r.fiberApp.Get("/torrents/:infohash/files/:filehash/streams/:name/watch", r.WatchForm)
	r.fiberApp.Get("/unpack/*", r.FileForm)
	// r.fiberApp.Get("/torrents/:infohash/files/:filehash/streams/:name/hls/*", r.HLSForm)

	return &r, nil
}

func errHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if code >= http.StatusInternalServerError {
		zerolog.Ctx(c.Context()).
			Error().
			Stack().Err(err).
			Int("code", code).
			Msg("http.internal.error")
	} else {
		zerolog.Ctx(c.Context()).
			Warn().
			Err(err).
			Int("code", code).
			Msg("http.client.error")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(code).SendString(err.Error())
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
		ctx = logger_utils.WithStrContextLog(ctx, "user_agent", c.Get("User-Agent"))
		c.SetContext(ctx)

		t0 := time.Now()

		err := c.Next()

		zerolog.Ctx(ctx).
			Debug().
			Str("elapsed", time.Since(t0).String()).
			// Str(c.Response().).
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

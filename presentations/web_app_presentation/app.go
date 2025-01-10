package web_app_presentation

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/pkg/errors"
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
		Views: html.NewFileSystem(http.FS(views.Static), ".html"),
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
		ctx = logger_utils.WithStrContextLog(ctx, "app_method", fmt.Sprintf("%s %s", c.Method(), c.Path()))
		c.SetContext(ctx)

		return c.Next()
	}
}

func (r *Presentation) Close(ctx context.Context) error {
	err := r.fiberApp.ShutdownWithContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error shutting down app")
	}

	return nil
}

//func (r *Presentation) torrentGetById(c fiber.Ctx) error {
//	idStr := c.Params("id")
//	id, err := uuid.Parse(idStr)
//	if err != nil {
//		// Send 400 on err
//		return errors.Wrap(err, "failed to parse id")
//	}
//
//	ctx := logger_utils.AddLoggerToCtx(c.Context())
//
//	torrent, err := r.torrentorService.GetTorrentMetadataByID(ctx, id)
//	if err != nil {
//		// Send 400 on err
//		return errors.Wrap(err, "failed to get torrent metadata")
//	}
//
//	// Send a string response to the client
//	return c.JSON(torrent)
//}

func (r *Presentation) Run(_ context.Context) error {
	err := r.fiberApp.Listen(settings.Settings.WebServer.URL)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	return nil
}

package web_app_presentation

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/pkg/errors"
	"net/http"
	"torrentor/presentations/web_app_presentation/views"
	"torrentor/services/torrentor_service_viewer"
	"torrentor/settings"
)

type Presentation struct {
	torrentorViewerService *torrentor_service_viewer.Service
	fiberApp               *fiber.App
}

func NewPresentation(
	ctx context.Context,
	torrentorViewerService *torrentor_service_viewer.Service,
) (*Presentation, error) {
	app := fiber.New(fiber.Config{
		Views: html.NewFileSystem(http.FS(views.Static), ".html"),
	})

	r := Presentation{
		torrentorViewerService: torrentorViewerService,
		fiberApp:               app,
	}

	// TODO move path to settings
	//r.fiberApp.Get("/api/torrent/:id", r.torrentGetById)
	r.fiberApp.Get("/", IndexForm)
	r.fiberApp.Get("/torrent/:id", r.TorrentForm)

	return &r, nil
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
//	torrent, err := r.torrentorViewerService.GetTorrentMetadataByID(ctx, id)
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

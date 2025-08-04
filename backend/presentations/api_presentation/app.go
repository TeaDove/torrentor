package api_presentation

import (
	"time"
	"torrentor/backend/services/torrentor_service"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/teadove/teasutils/fiber_utils"
)

type Presentation struct {
	torrentorService *torrentor_service.Service
}

func NewPresentation(torrentorService *torrentor_service.Service) *Presentation {
	return &Presentation{torrentorService: torrentorService}
}

func (r *Presentation) BuildApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ErrorHandler: fiber_utils.ErrHandler(),
	})
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Use(fiber_utils.MiddlewareLogger())
	app.Use(fiber_utils.MiddlewareCtxTimeout(10 * time.Second)) //nolint: mnd // don't care
	app.Use(cors.New(cors.ConfigDefault))

	api := app.Group("/api")
	api.Get("/torrents", r.listTorrents)
	api.Get("/stats", r.statsTorrents)

	api.Get("/torrents/:infoHash", r.getTorrent)
	api.Delete("/torrents/:infoHash", r.deleteTorrent)
	api.Post("/torrents/download", r.download)

	return app
}

func (r *Presentation) statsTorrents(c *fiber.Ctx) error {
	serviceStats, _, err := r.torrentorService.Stats(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to get stats")
	}

	return c.JSON(fiber.Map{"serviceStats": serviceStats})
}

func (r *Presentation) download(c *fiber.Ctx) error {
	type Request struct {
		Magnet string `json:"magnet" validate:"required"`
	}

	req, err := parseJSON[Request](c)
	if err != nil {
		return errors.WithStack(err)
	}

	torrent, err := r.torrentorService.DownloadAndSaveFromMagnet(c.UserContext(), req.Magnet)
	if err != nil {
		return errors.Wrap(err, "failed to start download")
	}

	return c.JSON(torrent)
}

func (r *Presentation) listTorrents(c *fiber.Ctx) error {
	torrents, err := r.torrentorService.ListOpenTorrents(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to get torrents")
	}

	return c.JSON(torrents)
}

func (r *Presentation) getTorrent(c *fiber.Ctx) error {
	infoHash := c.Params("infoHash")

	torrent, ok := r.torrentorService.GetTorrentByInfoHash(metainfo.NewHashFromHex(infoHash))
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "not found"}
	}

	return c.JSON(torrent)
}

func (r *Presentation) deleteTorrent(c *fiber.Ctx) error {
	infoHash := c.Params("infoHash")

	ok := r.torrentorService.DeleteTorrentByInfoHash(metainfo.NewHashFromHex(infoHash))
	if !ok {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "not found"}
	}

	return c.JSON(fiber.Map{"success": true})
}

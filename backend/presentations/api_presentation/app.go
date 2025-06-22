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
	api.Get("/torrents", r.getTorrents)
	api.Get("/torrents/:infoHash", r.getTorrent)
	api.Post("/torrents/download", r.download)

	return app
}

func (r *Presentation) download(c *fiber.Ctx) error {
	type Request struct {
		Magnet string `json:"magnet"`
	}

	var request Request

	err := c.BodyParser(&request)
	if err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.SendString(errors.Wrap(err, "unprocessable entity").Error())
	}

	torrent, _, err := r.torrentorService.DownloadAndSaveFromMagnet(c.UserContext(), request.Magnet)
	if err != nil {
		return errors.Wrap(err, "failed to start download")
	}

	return c.JSON(torrent)
}

func (r *Presentation) getTorrents(c *fiber.Ctx) error {
	torrents, err := r.torrentorService.GetAllTorrents(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to get torrents")
	}

	return c.JSON(torrents)
}

func (r *Presentation) getTorrent(c *fiber.Ctx) error {
	infoHash := c.Params("infoHash")

	torrent, err := r.torrentorService.GetTorrentByInfoHash(c.UserContext(), metainfo.NewHashFromHex(infoHash))
	if err != nil {
		return errors.Wrap(err, "failed to get torrent")
	}

	return c.JSON(torrent)
}

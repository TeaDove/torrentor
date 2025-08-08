package api_presentation

import (
	"time"
	"torrentor/backend/services/torrentor_service"
	"torrentor/backend/utils/validators_utils"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
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
	app.Use(fiber_utils.MiddlewareCtxTimeout(time.Minute))
	app.Use(cors.New(cors.ConfigDefault))

	api := app.Group("/api")
	api.Get("/torrents", r.listTorrents)
	api.Get("/stats", r.statsTorrents)

	api.Get("/torrents/:infoHash", r.getTorrent)
	api.Get("/torrents/:infoHash/files/:pathHash", r.getTorrentFile)
	api.Delete("/torrents/:infoHash", r.deleteTorrent)
	api.Post("/torrents/download", r.download)

	return app
}

func parseJSON[T any](c *fiber.Ctx) (T, error) {
	var v = new(T)

	err := c.BodyParser(v)
	if err != nil {
		return *v, sentUnprocessable(errors.Wrap(err, "failed to body parse"))
	}

	err = validators_utils.Validate.Struct(v)
	if err != nil {
		return *v, sentUnprocessable(errors.Wrap(err, "failed to validate"))
	}

	return *v, nil
}

func sentUnprocessable(err error) error {
	return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
}

func getInfoHashParams(c *fiber.Ctx) (metainfo.Hash, error) {
	var hash metainfo.Hash

	err := hash.FromHexString(c.Params("infoHash"))
	if err != nil {
		return metainfo.Hash{}, sentUnprocessable(errors.Wrap(err, "failed to get info hash"))
	}

	return hash, nil
}

package api_presentation

import (
	"torrentor/backend/services/torrentor_service"

	"github.com/gofiber/fiber/v2"
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
	app.Use(fiber_utils.MiddlewareLogger())
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("hi!") })

	return app
}

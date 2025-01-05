package web_app_presentation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

func RunApp() error {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	err := app.Listen(":3000")
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	return nil
}

package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := SetupApp()
	log.Fatal(app.Listen(":8000"))
}

func SetupApp() *fiber.App {
	app := fiber.New(config())

	app.Post("/converter", converter)
	app.Post("/resizer", resizer)
	app.Post("/compressor", compressor)

	return app
}

func config() fiber.Config {
	return fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = c.Status(code).JSON(map[string]any{
				"code": code,
				"message": e.Message,
				"status": "error",
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{
					"code": fiber.StatusInternalServerError,
					"message": "Internal Server Error",
					"status": "error",
				})
			}

			return nil
		},
	}
}

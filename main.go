package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Status(http.StatusOK).JSON(&fiber.Map{"msg": "app running succesfully"})
		return nil
	})

	app.Listen(":8000")

}

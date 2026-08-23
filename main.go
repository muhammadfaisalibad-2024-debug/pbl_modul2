package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students")
	})

	fmt.Println("Server berjalan di http://localhost:3000")

	app.Listen(":3000")
}

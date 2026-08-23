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

	api := app.Group("/api/v1")

	api.Get("/students", getStudents)
	api.Get("/students/:id", getStudent)
	api.Post("/students", createStudent)

	fmt.Println("Server berjalan di http://localhost:3000")

	app.Listen(":3000")
}

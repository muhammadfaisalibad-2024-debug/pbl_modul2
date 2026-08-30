package main

import (
	"fmt"
	"log"
	"strings"

	"api-students/config"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	cfg := config.Load()

	var err error
	db, err = config.NewDBPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost ||
			c.Method() == fiber.MethodPut ||
			c.Method() == fiber.MethodPatch {

			contentType := c.Get("Content-Type")

			if !strings.HasPrefix(contentType, fiber.MIMEApplicationJSON) {
				return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
					"success": false,
					"message": "Content-Type must be application/json",
				})
			}
		}

		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students")
	})

	api := app.Group("/api/v1")

	api.Get("/students", getStudents)
	api.Get("/students/:id", getStudent)
	api.Post("/students", createStudent)
	api.Put("/students/:id", updateStudent)
	api.Patch("/students/:id", patchStudent)
	api.Delete("/students/:id", deleteStudent)

	fmt.Println("Server berjalan di http://localhost:3000")

	app.Listen(":3000")
}

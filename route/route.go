package route

import (
	"api-students/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, svc *service.StudentService) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"message": "Database connection failed",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Database connection is healthy",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students")
	})

	api := app.Group("/api/v1")

	api.Get("/students", svc.GetStudents)
	api.Get("/students/:id", svc.GetStudent)
	api.Post("/students", svc.CreateStudent)
	api.Put("/students/:id", svc.UpdateStudent)
	api.Patch("/students/:id", svc.PatchStudent)
	api.Delete("/students/:id", svc.DeleteStudent)
}

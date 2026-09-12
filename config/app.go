package config

import (
	"log/slog"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
	"api-students/route"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApp(
	db *pgxpool.Pool,
	studentService *service.StudentService,
	authService *service.AuthService,
	jwtManager *helper.JWTManager,
	logger *slog.Logger,
) *fiber.App {

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024, // 1 MB

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	middleware.SetupMiddlewares(app, logger)

	route.RegisterRoutes(
		app,
		db,
		studentService,
		authService,
		jwtManager,
	)

	return app
}

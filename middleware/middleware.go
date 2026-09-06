package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupMiddlewares(app *fiber.App, logger *slog.Logger) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(cors.New())

	app.Use(RequestLogger(logger))
	app.Use(RequireJSON())
}

func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		reqID := c.Locals("requestid")
		if reqID == nil {
			reqID = ""
		}

		logger.Info("Incoming Request",
			slog.String("request_id", reqID.(string)),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.String("duration", duration.String()),
			slog.String("ip", c.IP()),
		)

		return err
	}
}

func RequireJSON() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

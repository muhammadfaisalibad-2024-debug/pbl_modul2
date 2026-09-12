package helper

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context(), 5*time.Second)
}

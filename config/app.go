package config

import (
	"errors"
	"log/slog"

	"api-students/app/model"
	"api-students/helper"
	"api-students/middleware"
	"api-students/route"

	"github.com/gofiber/fiber/v2"
)

func NewApp(
	logger *slog.Logger,
	deps route.Dependencies,
) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit:    1 * 1024 * 1024, // 1 MB
		ErrorHandler: newErrorHandler(logger),
	})

	middleware.SetupMiddlewares(app, logger)
	route.RegisterRoutes(app, deps)
	app.Use(func(c *fiber.Ctx) error { return helper.NotFound("endpoint tidak ditemukan") })

	return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *helper.AppError
		if !errors.As(err, &appErr) {
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				appErr = helper.BadRequest(fiberErr.Message)
			} else {
				appErr = helper.Internal(err)
			}
		}
		requestID, _ := c.Locals("requestid").(string)
		if appErr.Status >= 500 {
			logger.Error("request_failed", slog.String("request_id", requestID), slog.String("code", appErr.Code), slog.Int("status", appErr.Status), slog.String("error", appErr.Error()))
		} else {
			logger.Warn("request_rejected", slog.String("request_id", requestID), slog.String("code", appErr.Code), slog.Int("status", appErr.Status))
		}
		resp := model.ErrorResponse{
			Success:   false,
			Message:   appErr.Message,
			Errors:    appErr.Errors,
			RequestID: requestID,
		}
		if resp.Errors == nil && len(appErr.Fields) > 0 {
			resp.Errors = make(map[string][]string)
			for k, v := range appErr.Fields {
				resp.Errors[k] = []string{v}
			}
		}
		return c.Status(appErr.Status).JSON(resp)
	}
}

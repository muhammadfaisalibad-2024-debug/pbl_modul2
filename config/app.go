package config

import (
	"errors"
	"log/slog"

	"api-students/app/model"
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
	prestasiService *service.PrestasiService,
	authService *service.AuthService,
	userService *service.UserService,
	jwtManager *helper.JWTManager,
	permissions *helper.PermissionSet,
	logger *slog.Logger,
) *fiber.App {

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024, // 1 MB

		ErrorHandler: newErrorHandler(logger),
	})

	middleware.SetupMiddlewares(app, logger)

	route.RegisterRoutes(
		app,
		db,
		studentService,
		prestasiService,
		authService,
		userService,
		jwtManager,
		permissions,
	)
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
		return c.Status(appErr.Status).JSON(model.ErrorResponse{Success: false, Code: appErr.Code, Message: appErr.Message, Fields: appErr.Fields, RequestID: requestID})
	}
}

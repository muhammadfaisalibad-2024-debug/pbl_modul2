package route

import (
	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(
	app *fiber.App,
	db *pgxpool.Pool,
	studentService *service.StudentService,
	authService *service.AuthService,
	userService *service.UserService,
	jwtManager *helper.JWTManager,
	perms *helper.PermissionSet,
) {
	// Health check - public
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

	// Root endpoint - public
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students")
	})

	api := app.Group("/api/v1")

	// =========================
	// AUTHENTICATION
	// =========================
	auth := api.Group("/auth")

	auth.Post(
		"/register",
		middleware.RequireJSON(),
		authService.Register,
	)

	auth.Post(
		"/login",
		middleware.RequireJSON(),
		middleware.LoginRateLimiter(),
		authService.Login,
	)

	auth.Post(
		"/refresh",
		middleware.RequireJSON(),
		authService.Refresh,
	)

	auth.Post(
		"/logout",
		middleware.RequireJSON(),
		authService.Logout,
	)

	auth.Get(
		"/me",
		middleware.RequireAuth(jwtManager),
		authService.Me,
	)

	users := api.Group(
		"/users",
		middleware.RequireAuth(jwtManager),
	)

	users.Get("/", middleware.RequirePermission(perms, "user:list"), userService.List)
	users.Post("/", middleware.RequireJSON(), middleware.RequirePermission(perms, "user:update:any"), userService.Create)
	users.Get("/:id", userService.Get)
	users.Put("/:id", middleware.RequireJSON(), userService.Replace)
	users.Patch("/:id", middleware.RequireJSON(), userService.Patch)
	users.Delete("/:id", middleware.RequirePermission(perms, "user:delete"), userService.Delete)
	users.Patch("/:id/role", middleware.RequireJSON(), middleware.RequirePermission(perms, "role:assign"), userService.AssignRole)

	// =========================
	// STUDENTS - PROTECTED
	// =========================
	students := api.Group(
		"/students",
		middleware.RequireAuth(jwtManager),
	)

	students.Get(
		"/",
		middleware.RequirePermission(perms, "student:list"),
		studentService.GetStudents,
	)

	students.Get(
		"/:id",
		studentService.GetStudent,
	)

	students.Post(
		"/",
		middleware.RequireJSON(),
		middleware.RequirePermission(perms, "student:create"),
		studentService.CreateStudent,
	)

	students.Put(
		"/:id",
		middleware.RequireJSON(),
		studentService.UpdateStudent,
	)

	students.Patch(
		"/:id",
		middleware.RequireJSON(),
		studentService.PatchStudent,
	)

	students.Delete(
		"/:id",
		middleware.RequirePermission(perms, "student:delete"),
		studentService.DeleteStudent,
	)
}

package route

import (
	"api-students/app/model"
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
	courseService *service.CourseService,
	enrollmentService *service.EnrollmentService,
	authService *service.AuthService,
	jwtManager *helper.JWTManager,
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
		return c.JSON(fiber.Map{
			"success": true,
			"message": "SIAKAD Mini RESTful API Back End",
			"version": "1.0.0",
		})
	})

	api := app.Group("/api/v1")

	// ==========================================
	// 1. AUTHENTICATION ENDPOINTS
	// ==========================================
	auth := api.Group("/auth")

	// Endpoint 1: POST /api/v1/auth/login (Publik, Rate Limited 5x/menit)
	auth.Post(
		"/login",
		middleware.RequireJSON(),
		middleware.LoginRateLimiter(),
		authService.Login,
	)

	// Endpoint 2: GET /api/v1/auth/me (Semua Role)
	auth.Get(
		"/me",
		middleware.RequireAuth(jwtManager),
		authService.Me,
	)

	auth.Post("/refresh", middleware.RequireJSON(), authService.Refresh)
	auth.Post("/logout", authService.Logout)

	// ==========================================
	// 2. STUDENTS ENDPOINTS (Admin & Mahasiswa)
	// ==========================================
	students := api.Group(
		"/students",
		middleware.RequireAuth(jwtManager),
	)

	// Endpoint 3: GET /api/v1/students (Admin only)
	students.Get(
		"/",
		middleware.RequireRole(model.RoleAdmin),
		studentService.GetStudents,
	)

	// Endpoint 4: POST /api/v1/students (Admin only)
	students.Post(
		"/",
		middleware.RequireJSON(),
		middleware.RequireRole(model.RoleAdmin),
		studentService.CreateStudent,
	)

	// Endpoint 5: GET /api/v1/students/:id (Admin & Mahasiswa Data Sendiri)
	students.Get(
		"/:id",
		studentService.GetStudent,
	)

	// Endpoint 6: PUT /api/v1/students/:id (Admin only)
	students.Put(
		"/:id",
		middleware.RequireJSON(),
		middleware.RequireRole(model.RoleAdmin),
		studentService.UpdateStudent,
	)

	// Endpoint 7: DELETE /api/v1/students/:id (Admin only - Soft Delete)
	students.Delete(
		"/:id",
		middleware.RequireRole(model.RoleAdmin),
		studentService.DeleteStudent,
	)

	// ==========================================
	// 3. COURSES ENDPOINTS (Semua Role)
	// ==========================================
	courses := api.Group(
		"/courses",
		middleware.RequireAuth(jwtManager),
	)

	// Endpoint 8: GET /api/v1/courses (Semua Role)
	courses.Get(
		"/",
		courseService.GetCourses,
	)

	// ==========================================
	// 4. ENROLLMENTS (KRS) ENDPOINTS (Mahasiswa)
	// ==========================================
	enrollments := api.Group(
		"/enrollments",
		middleware.RequireAuth(jwtManager),
	)

	// Endpoint 9: POST /api/v1/enrollments (Mahasiswa only)
	enrollments.Post(
		"/",
		middleware.RequireJSON(),
		middleware.RequireRole(model.RoleMahasiswa),
		enrollmentService.CreateEnrollment,
	)

	// Endpoint 10: DELETE /api/v1/enrollments/:id (Mahasiswa milik sendiri)
	enrollments.Delete(
		"/:id",
		middleware.RequireRole(model.RoleMahasiswa),
		enrollmentService.DeleteEnrollment,
	)
}

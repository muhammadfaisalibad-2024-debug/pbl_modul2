package route

import (
	"context"
	"time"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Dependencies struct {
	Pool              *pgxpool.Pool
	Mongo             *mongo.Client
	JWT               *helper.JWTManager
	Permissions       *helper.PermissionSet
	StudentService    *service.StudentService
	CourseService     *service.CourseService
	EnrollmentService *service.EnrollmentService
	AuthService       *service.AuthService
	SupervisorService *service.SupervisorService
	ReviewService     *service.ReviewService
	PortfolioService  *service.PortfolioService
}

func RegisterRoutes(app *fiber.App, deps Dependencies) {
	// Health check - memeriksa DUA database (PostgreSQL & MongoDB)
	app.Get("/health", healthCheck(deps.Pool, deps.Mongo))
	app.Get("/api/v1/health", healthCheck(deps.Pool, deps.Mongo))

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
	if deps.AuthService != nil {
		auth := api.Group("/auth")
		auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
		auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)
		auth.Post("/refresh", deps.AuthService.Refresh)
		auth.Post("/logout", deps.AuthService.Logout)
	}

	// ==========================================
	// 2. STUDENTS & PORTFOLIO ENDPOINTS
	// ==========================================
	students := api.Group("/students", middleware.RequireAuth(deps.JWT))

	if deps.StudentService != nil {
		students.Get("/", middleware.RequireRole(model.RoleAdmin), deps.StudentService.GetStudents)
		students.Post("/", middleware.RequireRole(model.RoleAdmin), deps.StudentService.CreateStudent)
		students.Get("/:id", deps.StudentService.GetStudent)
		students.Put("/:id", middleware.RequireRole(model.RoleAdmin), deps.StudentService.UpdateStudent)
		students.Delete("/:id", middleware.RequireRole(model.RoleAdmin), deps.StudentService.DeleteStudent)
	}

	// Portofolio Mahasiswa (Modul 10)
	if deps.PortfolioService != nil {
		students.Get("/:id/portfolio", deps.PortfolioService.Get)
		students.Put("/:id/portfolio", deps.PortfolioService.Put)
		students.Post("/:id/portfolio/projects", deps.PortfolioService.AddProject)
		students.Delete("/:id/portfolio/projects/:projectId", deps.PortfolioService.RemoveProject)

		portfolios := api.Group("/portfolios", middleware.RequireAuth(deps.JWT))
		portfolios.Get("/stats/technologies", deps.PortfolioService.TopTechnologies)
	}

	// ==========================================
	// 3. COURSES, KRS & REVIEWS ENDPOINTS
	// ==========================================
	if deps.CourseService != nil {
		courses := api.Group("/courses", middleware.RequireAuth(deps.JWT))

		// "/mine" didaftarkan SEBELUM "/:id"
		courses.Get("/", deps.CourseService.List)
		courses.Get("/mine", deps.CourseService.Mine)
		courses.Get("/:id", deps.CourseService.Get)
		courses.Post("/:id/enroll", deps.CourseService.Enroll)
		courses.Delete("/:id/enroll", deps.CourseService.Drop)
		courses.Patch("/:id", middleware.RequirePermission(deps.Permissions, "course:manage"), deps.CourseService.UpdateQuota)

		// Ulasan Kelas (Modul 10 di dalam courses)
		if deps.ReviewService != nil {
			courses.Get("/:id/reviews", deps.ReviewService.ListByCourse)
			courses.Get("/:id/reviews/summary", deps.ReviewService.Summary)
			courses.Post("/:id/reviews", deps.ReviewService.Create)
		}
	}

	// ==========================================
	// 4. ENROLLMENTS (KRS) & SWAP
	// ==========================================
	enrollments := api.Group("/enrollments", middleware.RequireAuth(deps.JWT))
	if deps.CourseService != nil {
		enrollments.Post("/swap", deps.CourseService.Swap)
	}
	if deps.EnrollmentService != nil {
		enrollments.Post("/", middleware.RequireRole(model.RoleMahasiswa), deps.EnrollmentService.CreateEnrollment)
		enrollments.Delete("/:id", middleware.RequireRole(model.RoleMahasiswa), deps.EnrollmentService.DeleteEnrollment)
	}

	// ==========================================
	// 5. REVIEWS ENDPOINTS (Modul 10)
	// ==========================================
	if deps.ReviewService != nil {
		reviews := api.Group("/reviews", middleware.RequireAuth(deps.JWT))
		reviews.Get("/:id", deps.ReviewService.Get)
		reviews.Patch("/:id", deps.ReviewService.Update)
		reviews.Delete("/:id", deps.ReviewService.Delete)
		reviews.Post("/:id/replies", deps.ReviewService.Reply)
		reviews.Post("/:id/helpful", deps.ReviewService.MarkHelpful)
	}

	// ==========================================
	// 6. SUPERVISORS & SUPERVISIONS (Modul 9 D.2)
	// ==========================================
	if deps.SupervisorService != nil {
		supervisors := api.Group("/supervisors", middleware.RequireAuth(deps.JWT))
		supervisors.Post("/:id/assign", middleware.RequirePermission(deps.Permissions, "supervision:manage"), deps.SupervisorService.Assign)
		supervisors.Patch("/:id", middleware.RequirePermission(deps.Permissions, "supervision:manage"), deps.SupervisorService.UpdateQuota)

		supervisions := api.Group("/supervisions", middleware.RequireAuth(deps.JWT))
		supervisions.Post("/switch", deps.SupervisorService.Switch)
	}
}

func healthCheck(pool *pgxpool.Pool, mongoClient *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if pool != nil {
			if err := pool.Ping(ctx); err != nil {
				return helper.ServiceUnavailable("PostgreSQL tidak dapat dihubungi")
			}
		}

		if mongoClient != nil {
			if err := mongoClient.Ping(ctx, nil); err != nil {
				return helper.ServiceUnavailable("MongoDB tidak dapat dihubungi")
			}
		}

		return helper.Success(c, fiber.StatusOK, "server, PostgreSQL, dan MongoDB berjalan", nil)
	}
}

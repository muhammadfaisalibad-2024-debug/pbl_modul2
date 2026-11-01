package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
	"api-students/route"
)

const minSecretLength = 32

func main() {
	config.LoadEnv()
	logger := config.NewLogger()

	// Satu program, tiga cara menjalankan:
	//   go run . -role=api    hanya melayani HTTP
	//   go run . -role=worker hanya relay outbox dan worker queue
	//   go run . -role=all    keduanya dalam satu proses (bawaan)
	role := flag.String("role", config.GetEnv("APP_ROLE", "all"), "api | worker | all")
	flag.Parse()

	runAPI := *role == "api" || *role == "all"
	runWorker := *role == "worker" || *role == "all"
	if !runAPI && !runWorker {
		logger.Error("role tidak dikenal", slog.String("role", *role))
		os.Exit(1)
	}

	cfg := config.Load()
	pool, err := database.NewDBPool(cfg)
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// ctx dibatalkan ketika Ctrl+C ditekan atau sistem mengirim SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var app *fiber.App
	if runAPI {
		app = startAPI(logger, pool, cfg)
		if app == nil {
			os.Exit(1)
		}
	}

	var stopWorker func()
	if runWorker {
		stopWorker = startWorker(ctx, logger, pool)
		if stopWorker == nil {
			os.Exit(1)
		}
	}

	logger.Info("aplikasi berjalan", slog.String("role", *role))
	<-ctx.Done()

	// Urutan berhenti: pintu masuk ditutup lebih dulu (HTTP), baru kemudian
	// pekerja di belakangnya.
	if app != nil {
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			logger.Error("shutdown HTTP tidak rapi", slog.String("error", err.Error()))
		}
	}

	if stopWorker != nil {
		stopWorker()
	}

	logger.Info("aplikasi berhenti dengan rapi")
}

// startAPI menyiapkan seluruh kebutuhan HTTP. Mengembalikan nil bila gagal.
func startAPI(logger *slog.Logger, pool *pgxpool.Pool, cfg config.Config) *fiber.App {
	jwtSecret := cfg.JWTSecret
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET tidak diisi atau terlalu pendek",
			slog.Int("minimal_karakter", minSecretLength))
		return nil
	}

	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke MongoDB", slog.String("error", err.Error()))
		return nil
	}
	mongoDB := mongoClient.Database(config.GetEnv("MONGO_DATABASE", "praktikum_backend"))

	accessMinutes, err := strconv.Atoi(cfg.JWTAccessTTLMinutes)
	if err != nil || accessMinutes <= 0 {
		accessMinutes = 15
	}

	refreshDays, err := strconv.Atoi(cfg.JWTRefreshTTLDays)
	if err != nil || refreshDays <= 0 {
		refreshDays = 7
	}

	jwtManager := helper.NewJWTManager(
		jwtSecret,
		cfg.JWTIssuer,
		time.Duration(accessMinutes)*time.Minute,
	)

	// Repositories
	userRepository := repository.NewUserRepository(pool)
	studentRepository := repository.NewStudentRepository(pool)
	tokenRepository := repository.NewTokenRepository(pool)
	courseRepository := repository.NewCourseRepository(pool)
	enrollmentRepository := repository.NewEnrollmentRepository(pool)
	supervisorRepository := repository.NewSupervisorRepository(pool)
	outboxRepository := repository.NewOutboxRepository(pool)
	notificationRepository := repository.NewNotificationRepository(pool)
	txManager := repository.NewTxManager(pool)

	reviewRepository := repository.NewReviewRepository(mongoDB)
	if err := reviewRepository.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index MongoDB", slog.String("error", err.Error()))
		return nil
	}

	portfolioRepository := repository.NewPortfolioRepository(mongoDB)
	if err := portfolioRepository.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index portfolio MongoDB", slog.String("error", err.Error()))
		return nil
	}

	permissions, err := userRepository.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("gagal memuat permission", slog.String("error", err.Error()))
		return nil
	}

	// Inspector hanya MEMBACA keadaan queue untuk endpoint admin.
	inspector := asynq.NewInspector(config.RedisOpt())

	app := config.NewApp(logger, route.Dependencies{
		Pool:              pool,
		Mongo:             mongoClient,
		JWT:               jwtManager,
		Permissions:       permissions,
		StudentService:    service.NewStudentService(studentRepository),
		CourseService:     service.NewCourseService(courseRepository, outboxRepository, txManager),
		EnrollmentService: service.NewEnrollmentService(enrollmentRepository, studentRepository),
		AuthService: service.NewAuthService(
			userRepository,
			studentRepository,
			tokenRepository,
			jwtManager,
			time.Duration(refreshDays)*24*time.Hour,
		),
		SupervisorService:   service.NewSupervisorService(supervisorRepository, txManager),
		ReviewService:       service.NewReviewService(reviewRepository, courseRepository, permissions),
		PortfolioService:    service.NewPortfolioService(portfolioRepository, studentRepository, permissions),
		NotificationService: service.NewNotificationService(notificationRepository, outboxRepository, inspector),
	})

	// Sumber daya milik API ditutup setelah server HTTP berhenti.
	app.Hooks().OnShutdown(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = inspector.Close()
		return mongoClient.Disconnect(ctx)
	})

	port := cfg.AppPort
	if port == "" {
		port = "3000"
	}

	go func() {
		fmt.Printf("SIAKAD Mini API Server berjalan di http://localhost:%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("HTTP berjalan", slog.String("port", port))
	return app
}

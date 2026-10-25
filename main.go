package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
	"api-students/route"
)

func main() {
	// 1. Load env
	cfg := config.Load()

	// 2. Create logger
	logger := config.NewLogger()

	// 3. Validate JWT secret
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET harus diisi minimal 32 karakter")
	}

	// 4. Create database pool (PostgreSQL)
	db, err := database.NewDBPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// 5. Connect to MongoDB
	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(ctx)
	}()
	mongoDB := mongoClient.Database(config.GetEnv("MONGO_DATABASE", "praktikum_backend"))

	// ==============================
	// JWT CONFIG
	// ==============================
	accessMinutes, err := strconv.Atoi(cfg.JWTAccessTTLMinutes)
	if err != nil || accessMinutes <= 0 {
		accessMinutes = 15
	}

	refreshDays, err := strconv.Atoi(cfg.JWTRefreshTTLDays)
	if err != nil || refreshDays <= 0 {
		refreshDays = 7
	}

	jwtManager := helper.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		time.Duration(accessMinutes)*time.Minute,
	)

	// ==============================
	// REPOSITORIES
	// ==============================
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	txManager := repository.NewTxManager(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	supervisorRepo := repository.NewSupervisorRepository(db)

	reviewRepo := repository.NewReviewRepository(mongoDB)
	if err := reviewRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index review MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	portfolioRepo := repository.NewPortfolioRepository(mongoDB)
	if err := portfolioRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index portfolio MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// ==============================
	// LOAD PERMISSIONS (sekali saat aplikasi menyala)
	// ==============================
	permissions, err := userRepo.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("gagal memuat permissions", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// ==============================
	// SERVICES
	// ==============================
	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo, txManager)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, studentRepo)
	authService := service.NewAuthService(
		userRepo,
		studentRepo,
		tokenRepo,
		jwtManager,
		time.Duration(refreshDays)*24*time.Hour,
	)
	supervisorService := service.NewSupervisorService(supervisorRepo, txManager)
	reviewService := service.NewReviewService(reviewRepo, courseRepo, permissions)
	portfolioService := service.NewPortfolioService(portfolioRepo, studentRepo, permissions)

	// ==============================
	// FIBER APP
	// ==============================
	app := config.NewApp(logger, route.Dependencies{
		Pool:              db,
		Mongo:             mongoClient,
		JWT:               jwtManager,
		Permissions:       permissions,
		StudentService:    studentService,
		CourseService:     courseService,
		EnrollmentService: enrollmentService,
		AuthService:       authService,
		SupervisorService: supervisorService,
		ReviewService:     reviewService,
		PortfolioService:  portfolioService,
	})

	// ==============================
	// RUN SERVER
	// ==============================
	port := cfg.AppPort
	if port == "" {
		port = "3000"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		fmt.Printf(
			"SIAKAD Mini API Server berjalan di http://localhost:%s\n",
			port,
		)

		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	<-quit

	log.Println("Shutting down server...")

	if err := app.ShutdownWithContext(
		context.Background(),
	); err != nil {
		log.Fatalf(
			"Server shutdown error: %v",
			err,
		)
	}

	log.Println("Server stopped.")
}

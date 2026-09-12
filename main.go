package main

import (
	"context"
	"fmt"
	"log"
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

	// 4. Create database pool
	db, err := database.NewDBPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

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

	studentRepo := repository.NewStudentRepository(db)
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// ==============================
	// SERVICES
	// ==============================

	studentService := service.NewStudentService(studentRepo)

	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		jwtManager,
		time.Duration(refreshDays)*24*time.Hour,
	)

	// ==============================
	// FIBER APP
	// ==============================

	app := config.NewApp(
		db,
		studentService,
		authService,
		jwtManager,
		logger,
	)

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
			"Server berjalan di http://localhost:%s\n",
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

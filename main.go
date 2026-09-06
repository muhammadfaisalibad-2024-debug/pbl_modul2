package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
)

func main() {
	// 1. Load env
	cfg := config.Load()

	// 2. Create logger
	logger := config.NewLogger()

	// 3. Create database pool
	db, err := database.NewDBPool(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// 4. Create repository
	repo := repository.NewStudentRepository(db)

	// 5. Create service
	svc := service.NewStudentService(repo)

	// 6. Create Fiber app
	app := config.NewApp(db, svc, logger)

	// 7. Run server (graceful shutdown)
	port := cfg.AppPort
	if port == "" {
		port = "3000"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server berjalan di http://localhost:%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 8. Graceful shutdown on signal
	<-quit
	log.Println("Shutting down server...")
	if err := app.ShutdownWithContext(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped.")
}

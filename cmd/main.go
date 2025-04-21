package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-workmate/internal/config"
	"test-workmate/internal/repository/database"
	"test-workmate/internal/server/handler"
	server "test-workmate/internal/server/http"
	"test-workmate/internal/service"
	"test-workmate/pkg/database/postgres"
	"test-workmate/pkg/logger"

	_ "github.com/lib/pq"
)

// @title           Test Workmate API
// @version         1.0
// @description     Простой сервис для задач

// @host      localhost:8080
// @BasePath  /

func main() {
	logger.ZapLoggerInit()
	ctx := context.Background()

	cfg := config.MustInit(os.Getenv("IS_PROD"))

	// Init Postgres
	mainPC, err := postgres.NewPostgresConnection(&cfg.PostgresMain)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to Postgres: %v", err))
	}
	defer mainPC.Close()

	logger.Info("postgres connection established!")

	taskRepo := database.NewTaskRepo(mainPC)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	router := server.InitRouter(taskHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting HTTP server on :%d", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	awaitStop(ctx, srv)
}

func awaitStop(ctx context.Context, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit
	logger.Info(fmt.Sprintf("Shutting down server... signal: %v", sig))

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("HTTP server Shutdown: %v", err))
	}
}

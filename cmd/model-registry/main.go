package main

import (
	"log/slog"
	"os"
	"strconv"

	"model-registry/internal/config"
	"model-registry/internal/handlers"
	"model-registry/internal/models"
	"model-registry/internal/repositories"
	"model-registry/internal/routes"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.GetConfig()

	logger := setupLogger(cfg.Env)

	logger.Info("starting app", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	if err := config.InitDB(cfg.DatabasePath, &models.Model{}, &models.ModelVersion{}); err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	db := config.GetDB()
	logger.Info("Database initialized successfully")

	modelRepo := repositories.NewModelRepository(db)
	versionRepo := repositories.NewVersionRepository(db)
	handlers.InitHandlers(modelRepo, versionRepo)
	logger.Info("Repositories initialized")

	router := routes.SetupRouter()
	logger.Info("Routes set up")

	port := strconv.Itoa(cfg.Port)
	logger.Info("Starting server", slog.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

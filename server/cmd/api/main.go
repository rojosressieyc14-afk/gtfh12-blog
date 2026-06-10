package main

import (
	"os"

	"blog/server/internal/config"
	"blog/server/internal/database"
	"blog/server/internal/logger"
	"blog/server/internal/router"
	"blog/server/internal/service"
	"blog/server/internal/utils"
	"go.uber.org/zap"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	logger.Init(
		getEnv("LOG_LEVEL", "info"),
		getEnv("LOG_FORMAT", "text"),
	)
	defer logger.Sync()

	l := logger.L

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		l.Fatal("config invalid", zap.Error(err))
	}
	if err := utils.InitKeys(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath); err != nil {
		l.Fatal("init JWT keys", zap.Error(err))
	}
	l.Info("server starting", zap.String("config", cfg.StartupSummary()))

	db, err := database.Connect(cfg)
	if err != nil {
		l.Fatal("connect database", zap.Error(err))
	}

	if err := database.Migrate(db); err != nil {
		l.Fatal("migrate database", zap.Error(err))
	}

	if err := database.SeedAdmin(db, cfg); err != nil {
		l.Fatal("seed admin", zap.Error(err))
	}
	if err := database.SeedCategories(db); err != nil {
		l.Fatal("seed categories", zap.Error(err))
	}
	if err := service.LoadModerationSettings(db); err != nil {
		l.Fatal("load moderation settings", zap.Error(err))
	}
	if err := service.LoadSensitiveWords(db); err != nil {
		l.Fatal("load sensitive words", zap.Error(err))
	}

	r := router.New(cfg, db)
	l.Info("server listening", zap.String("port", cfg.ServerPort))
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		l.Fatal("run server", zap.Error(err))
	}
}

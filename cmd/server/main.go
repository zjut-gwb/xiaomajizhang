package main

import (
	"fmt"
	"log"

	"github.com/asikeida/xiaomajizhang/internal/config"
	"github.com/asikeida/xiaomajizhang/internal/database"
	"github.com/asikeida/xiaomajizhang/internal/handler"
	"github.com/asikeida/xiaomajizhang/internal/logger"
	"github.com/asikeida/xiaomajizhang/internal/repository"
	"github.com/asikeida/xiaomajizhang/internal/router"
	"github.com/asikeida/xiaomajizhang/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	zapLogger, err := logger.New(cfg.App.Env, cfg.Log.Level)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()

	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		zapLogger.Fatal("connect database failed", zap.Error(err))
	}
	if err := database.AutoMigrate(db); err != nil {
		zapLogger.Fatal("auto migrate failed", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := service.NewAuthService(cfg.JWT, userRepo, refreshTokenRepo, zapLogger)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(authService, zapLogger)
	userHandler := handler.NewUserHandler(userService, zapLogger)

	r := router.New(cfg.JWT.Secret, router.Handlers{
		Auth: authHandler,
		User: userHandler,
	})

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	zapLogger.Info("server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		zapLogger.Fatal("server stopped", zap.Error(err))
	}
}

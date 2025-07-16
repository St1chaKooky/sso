package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//иницилизация конфига

	cfg := config.MustLoadConfig()

	//иницилизация логера

	logger := setupLogger()
	defer func(logger *slog.Logger) {

	}(logger)

	logger.Info("start app", "env", cfg.Env)

	//иницилизация приложения

	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	//запуск grpc сервера

	go application.GRPCSrv.MustRun()

	//Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	logger.Info("application stopped", "signal", sign.String())

	application.GRPCSrv.Stop()
}

func setupLogger() *slog.Logger {
	core := slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, &core)
	logger := slog.New(handler)
	return logger
}

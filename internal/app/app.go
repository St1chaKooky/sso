package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storesPath string,
	tokenTTL time.Duration,
) *App {
	//иницилизировать хранилище
	storage, err := sqlite.NewStorage(storesPath)
	if err != nil {
		panic(err)
	}
	//иницилизировать auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}

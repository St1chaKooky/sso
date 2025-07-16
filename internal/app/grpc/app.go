package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() error {
	if err := a.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		"op", op,
		"port", a.port,
	)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%w: %s", err, op)
	}
	log.Info("grpc server is running", "addr", l.Addr().String())
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%w: %s", err, op)
	}
	return nil
}

func (a *App) Stop() error {
	const op = "grpcapp.Stop"
	a.log.With("op", op).Info("stopping grpc server", "port", a.port)

	a.gRPCServer.GracefulStop() // завершение разных операций и ток потом выключаем сервер
	return nil
}

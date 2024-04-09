package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/config"
	authgrpc "github.com/rautaruukkipalich/go_auth_grpc/internal/grpc/auth"
	authsrvcs "github.com/rautaruukkipalich/go_auth_grpc/internal/services/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.ConnectionTimeout(
			cfg.Server.ConnTimeout,
		),
	)

	auth := authsrvcs.New(
		nil,
		nil, 
		nil,
		nil,
		log,
		cfg.Token.TTL,
	)

	authgrpc.RegisterServer(gRPCServer, auth)

	return &App{
		log: log,
		gRPCServer: gRPCServer,
		port: cfg.Server.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("run grpc server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("stop grpc server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/app/kafka"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/config"
	authgrpc "github.com/rautaruukkipalich/go_auth_grpc/internal/grpc/auth"
	authsrvcs "github.com/rautaruukkipalich/go_auth_grpc/internal/services/auth"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/storage/sqlstorage"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	broker     kafka.Brokerer
	port       string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	broker kafka.Brokerer,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.ConnectionTimeout(
			cfg.Server.ConnTimeout,
		),
	)

	var dbURI string

	switch cfg.Database.Driver {
	case "postgres":
		dbURI = fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.Driver,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)
	default:
		panic("invalid database driver")
	}

	storage, err := sqlstorage.New(dbURI)
	if err != nil {
		panic(err)
	}

	auth := authsrvcs.New(
		storage,
		storage,
		storage,
		storage,
		log,
		cfg.Token.TTL,
		broker,
	)

	authgrpc.RegisterServer(gRPCServer, auth)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.Server.Port,
		broker:     broker,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.grpc.app.Run"
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
	const op = "app.grpc.app.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("stop grpc server", slog.String("port", a.port))

	a.broker.Stop()
	a.gRPCServer.GracefulStop()
}

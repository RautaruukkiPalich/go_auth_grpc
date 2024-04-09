package app

import (
	"log/slog"

	grpcapp "github.com/rautaruukkipalich/go_auth_grpc/internal/app/grpc"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/config"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// init storage 
	// init service auth

	grpcApp := grpcapp.New(log, cfg)

	return &App{
		GRPCSrv: grpcApp,
	}
}

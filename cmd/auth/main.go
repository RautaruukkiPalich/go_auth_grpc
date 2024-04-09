package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/app"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/config"
	// "github.com/rautaruukkipalich/go_auth_grpc/internal/lib/slerr"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()

	log := MustRunLogger(cfg.Env)

	log.Info("logger initialized")
	if cfg.Env == envLocal {
		log.Info("", slog.Any("config", cfg))
	}

	application := app.New(log, cfg)
	go application.GRPCSrv.MustRun()

	// TODO: init app

	// TODO: run server

	// stop application
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<- stop
	
	application.GRPCSrv.Stop()
	log.Info("application stopped")

}

func MustRunLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("unknown env: " + env)
	}
	return log
}
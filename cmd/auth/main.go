package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/app"
	"github.com/rautaruukkipalich/go_auth_grpc/internal/config"
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

	// init app
	application := app.New(log, cfg)

	// run server
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<- stop
	
	// stop server
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

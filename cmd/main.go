package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad("./config/local.yaml")

	log := setupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.Storage.Type, cfg.TokenTTL)

	go application.GRPCServ.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop

	log.Info("Stopping application", slog.String("stop", sig.String()))

	application.GRPCServ.Stop()

	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		file, err := os.OpenFile("./logs/log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0224)
		if err != nil {
			slog.Error("Failed to open log file", "error", err)
			return nil
		}

		log = slog.New(
			slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

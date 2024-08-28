package app

import (
	grpcapp "auth/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServ *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storage string, tokenTTL time.Duration) *App {

	// TODO: asd

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServ: grpcApp,
	}
}

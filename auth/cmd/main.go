package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jekiti/citydrive/auth/internal/app"
	"github.com/jekiti/citydrive/pkg/logger"
)

func main() {
	log := logger.SetupLogger("DEBUG")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.NewApp(log, "./.env")
	if err != nil {
		log.Error("failed to create app:", slog.Any("error", err))
		return
	}
	defer app.Close()

	log.Info("starting app...")
	if err := app.Run(ctx); err != nil {
		log.Error("failed to run app:", slog.Any("error", err))
		return
	}
	log.Info("app stopped")
}


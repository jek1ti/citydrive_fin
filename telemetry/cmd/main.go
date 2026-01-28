package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jekiti/citydrive/telemetry/internal/app"
	"github.com/jekiti/citydrive/telemetry/internal/config"
)

func main() {
	cfg := config.LoadTelemetryConfig()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.NewApp(cfg, log, ".env")
	if err != nil {
		log.Error("error creating app in main", "error", err)
		os.Exit(1)
	}
	err = app.Run(ctx)
	if err != nil {
		log.Error("app run failed", "error", err)
		os.Exit(1)
	}
	log.Info("application stopped")
}

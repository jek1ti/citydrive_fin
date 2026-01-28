package app

import (
	"context"
	"log/slog"

	telemetrypb "github.com/jekiti/citydrive/gen/proto/telemetry"
	"github.com/jekiti/citydrive/telemetry/internal/config"
	"github.com/jekiti/citydrive/telemetry/internal/handler"
	"github.com/jekiti/citydrive/telemetry/internal/producer"
	"github.com/jekiti/citydrive/telemetry/internal/repository"
	"github.com/jekiti/citydrive/telemetry/internal/server"
	"github.com/jekiti/citydrive/telemetry/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	log      *slog.Logger
	port     string
	register func(*grpc.Server)
}

func NewApp(cfg *config.TelemetryConfig, log *slog.Logger, envPath string) (*App, error) {
	log = log.With("module", "app", "function", "NewApp")
	redis, err := repository.NewRedisRepository(cfg, log)
	if err != nil {
		log.Error("error creating repo in app", "error",  err)
		return nil, err
	}

	violationService := service.NewViolationService(&cfg.Violations, log)
	producerKafka, err := producer.NewKafkaProducer(cfg, log)
	if err != nil {
		log.Error("error creating producer in app", "error", err)
		return nil, err
	}

	telemetryService := service.NewTelemetryService(redis, violationService, producerKafka, cfg, log)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService, log)
	reg := func(s *grpc.Server) {
		telemetrypb.RegisterTelemetryServiceServer(s, telemetryHandler)
	}
	log.Info("app initialized successfully")
	return &App{
		log:      log,
		port:     cfg.GRPC.Port,
		register: reg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	log := a.log.With("function", "Run")
	 log.Info("starting app")
	return server.Run(ctx, a.log, a.port, a.register)
}

func (a *App) Close() {
	a.log.Info("app closed")
}

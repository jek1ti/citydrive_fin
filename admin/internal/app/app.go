package app

import (
	"context"
	"log/slog"

	"github.com/jekiti/citydrive/admin/internal/config"
	"github.com/jekiti/citydrive/admin/internal/handlers"
	"github.com/jekiti/citydrive/admin/internal/repository"
	"github.com/jekiti/citydrive/admin/internal/server"
	"github.com/jekiti/citydrive/admin/internal/service"
	adminpb "github.com/jekiti/citydrive/gen/proto/admin"
	"google.golang.org/grpc"
)

type App struct {
	log      *slog.Logger
	port     string
	register func(*grpc.Server)
}

func NewApp(cfg *config.AdminConfig, log *slog.Logger) (*App, error) {
	log = log.With("module", "app", "function", "NewApp")
	redis, err := repository.NewRedisRepository(cfg, log)
	if err != nil{
		log.Error("error creating redis repo in app", "error", err)
		return nil, err
	}
	postgres, err := repository.NewPostgresRepository(&cfg.Postgres, log)
	if err != nil {
		log.Error("error creating postgres repo in app", "error", err)
		return nil, err
	}

	adminService := service.NewService(postgres, redis, log)

	adminHandler := handlers.NewHandler(adminService, log)
	reg := func(s *grpc.Server) {
		adminpb.RegisterAdminServiceServer(s, adminHandler)
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
	return server.Run(ctx, log, a.port, a.register)
}

func (a *App) Close() {
	a.log.Info("app closed")
}

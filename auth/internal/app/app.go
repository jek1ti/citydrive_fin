package app

import (
	"context"
	"log/slog"

	"github.com/jekiti/citydrive/auth/internal/config"
	"github.com/jekiti/citydrive/auth/internal/handler"
	authrepository "github.com/jekiti/citydrive/auth/internal/repository"
	"github.com/jekiti/citydrive/auth/internal/server"
	authservice "github.com/jekiti/citydrive/auth/internal/service"
	"github.com/jekiti/citydrive/auth/postgres"
	auth "github.com/jekiti/citydrive/gen/proto/auth"
	"google.golang.org/grpc"
)

type App struct {
	log      *slog.Logger
	db       *postgres.Postgres
	port     string
	register func(*grpc.Server)
}

func NewApp(log *slog.Logger, envPath string) (*App, error) {
	cfg := config.LoadConfig(envPath)
	db, err := postgres.NewPostgres(cfg.Database)
	if err != nil {
		log.Error("failed to connect to postgres:", slog.Any("error", err))
		return nil, err
	}

	pool := db.Master()
	repo := authrepository.NewUserRepository(pool, log)
	service := authservice.NewUserService(repo, log, cfg.JWT.SecretKey)
	authHandler := handler.NewAuthHandler(service, log)
	reg := func(s *grpc.Server) {
		auth.RegisterAuthServiceServer(s, authHandler)
	}
	return &App{
		log:      log,
		db:       db,
		port:     cfg.Server.GRPCPort,
		register: reg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return server.Run(ctx, a.log, a.port, a.register)
}

func (a *App) Close() {
	a.log.Info("app closed")
	if a.db != nil {
		a.log.Info("db connection closing...")
		a.db.Close()
		a.log.Info("db connection closed")
	}
}

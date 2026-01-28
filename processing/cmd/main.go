package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"log/slog"

	"github.com/jekiti/citydrive/processing/internal/app"
	"github.com/jekiti/citydrive/processing/internal/config"
	"github.com/jekiti/citydrive/processing/internal/repository"
	"github.com/jekiti/citydrive/processing/internal/service"
)

func main() {
	cfg := config.LoadProcessorConfig()

	router := app.NewServer()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	consumer := repository.NewKafkaConsumer(&cfg.Kafka, log)
	cache, err := repository.NewRedisRepository(&cfg.Redis, log)
	if err != nil {
		panic(err)
	}

	repo, err := repository.NewPostgresRepository(&cfg.DB, log)
	if err != nil {
		panic(err)
	}

	svc := service.NewService(consumer, cache, repo, &cfg.Processor, log)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := svc.ProcessTelemetry(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Error("process telemetry", "error", err)
		}
	}()

	addr := ":" + cfg.App.HTTPPort
	srv := &http.Server{
		Addr:              addr,
		Handler:           router, 
		ReadHeaderTimeout: 5 * time.Second,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("http listen failed", "addr", addr, "error", err)
		return
	}
	log.Info("HTTP server listening", "addr", addr)

	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http serve failed", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("Shutdown Server ...")
	

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown", "error", err)
	}
	wg.Wait()
	
	if err := repo.Close(); err != nil {
		log.Error("close postgres", "error", err)
	}
	if err := cache.Close(); err != nil {
		log.Error("close redis", "error", err)
	}
	if err := consumer.Close(); err != nil {
		log.Error("close kafka", "error", err)
	}

	log.Info("Shutdown completed")
}

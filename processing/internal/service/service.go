package service

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/jekiti/citydrive/processing/internal/config"
	"github.com/jekiti/citydrive/processing/internal/domain"
	"github.com/jekiti/citydrive/processing/internal/repository"
)

type Service interface {
	ProcessTelemetry(ctx context.Context) error
}

type ProcessingService struct {
	consumer   repository.Consumer
	cache      repository.CacheRepository
	log        *slog.Logger
	repository repository.DBRepository
	config     *config.ProcessorSpecificConfig
}

func NewService(consumer repository.Consumer, cache repository.CacheRepository, repo repository.DBRepository, config *config.ProcessorSpecificConfig, log *slog.Logger) Service {
	return &ProcessingService{
		consumer:   consumer,
		cache:      cache,
		repository: repo,
		log:        log,
		config:     config,
	}
}

func (s *ProcessingService) ProcessTelemetry(ctx context.Context) error {
	log := s.log.With("module", "service", "function", "ProcessTelemetry")
	log.Info("processing telemetry data")

	for {
		select {
		case <-ctx.Done():
			log.Info("shutting down telemetry processing")
			return nil
		default:
			log.Info("getting messages from service")
			messages, err := s.consumer.GetMessages(ctx, 1)
			if err != nil {
				log.Error("error getting messages from consumer", "error", err)
				continue
			}
			log.Info("fetched messages from consumer", "count", len(messages))

			for _, msg := range messages {
				prev, err := s.cache.GetCarState(msg.CarID)
				if err != nil {
					log.Error("error getting car state from cache", "error", err)
					continue
				}
				if hasDataChanged(prev, &msg) {
					err := s.cache.SaveCarState(msg)
					if err != nil {
						log.Error("error saving car state to cache", "error", err)
						continue
					}
				}
				err = s.repository.SaveTelemetry(msg)
				if err != nil {
					log.Error("error saving telemetry to repository", "error", err)
					continue
				}
			}
			if len(messages) != 0{
				log.Info("commiting messages", "msg", messages[0])
			}
			err = s.consumer.Commit()
			if err != nil {
				log.Error("error committing messages", "error", err)
				continue
			}
			if len(messages) == 0 {
				log.Info("no new messages to process, sleeping", "duration", s.config.PollTimeout)
				time.Sleep(s.config.PollTimeout)
			}
		}
	}

}

func hasDataChanged(previous, current *domain.CarTelemetry) bool {
	if previous == nil {
		return true
	}

	return previous.Brand != current.Brand ||
		previous.Model != current.Model ||
		previous.YearOfManufacture != current.YearOfManufacture ||
		previous.Odo != current.Odo ||
		!floatsEqual(previous.Lat, current.Lat) ||
		!floatsEqual(previous.Lon, current.Lon) ||
		!floatsEqual(previous.Fuel, current.Fuel) ||
		previous.FuelType != current.FuelType ||
		previous.Speed != current.Speed ||
		previous.EngineOn != current.EngineOn ||
		previous.Locked != current.Locked ||
		previous.Activated != current.Activated ||
		previous.Rpm != current.Rpm ||
		previous.Handbrake != current.Handbrake
}

func floatsEqual(a, b float64) bool {
	const epsilon = 0.000001
	return math.Abs(a-b) < epsilon
}

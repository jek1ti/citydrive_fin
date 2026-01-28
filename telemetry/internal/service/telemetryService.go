package service

import (
	"context"
	"log/slog"
	"math"

	"github.com/jekiti/citydrive/telemetry/internal/config"
	"github.com/jekiti/citydrive/telemetry/internal/models"
	"github.com/jekiti/citydrive/telemetry/internal/producer"
	"github.com/jekiti/citydrive/telemetry/internal/repository"
)

type TelemetryService struct {
	redis            *repository.RedisRepository
	violationService *ViolationService
	producer         *producer.KafkaProducer
	config           *config.TelemetryConfig
	log              *slog.Logger
}

func NewTelemetryService(redis *repository.RedisRepository,
	violationService *ViolationService,
	producer *producer.KafkaProducer,
	config *config.TelemetryConfig,
	log *slog.Logger) *TelemetryService {
	return &TelemetryService{
		redis:            redis,
		violationService: violationService,
		producer:         producer,
		config:           config,
		log:              log,
	}
}

func (s *TelemetryService) ProcessTelemetry(ctx context.Context, carID string, data *models.TelemetryData) error {
	traceID := ctx.Value("trace_id")
	log := s.log.With(
		"module", "telemetry.service",
		"function", "ProcessTelemetry",
		"car_id", carID,
		"trace_id", traceID,
	)

	log.Info("start processing telemetry")
	err := s.producer.SendTelemetry(ctx, carID, data)
	if err != nil {
		log.Error("error sending telemetry", "error", err)
		return err
	}
	log.Info("getting car state telemetry")
	prev, err := s.redis.GetCarState(ctx, carID)
	if err != nil {
		log.Error("error getting car state", "error", err)
		return err
	}
	log.Info("setting car state")
	if hasDataChanged(prev, data) {
		err := s.redis.SetCarState(ctx, carID, data)
		if err != nil {
			log.Error("error setting car state", "error", err)
			return err
		}

		violations := s.violationService.CheckViolations(ctx, carID, data)
		log.Info("violations detected", "count", len(violations))
		for _, violation := range violations {
			err = s.producer.SendViolation(ctx, violation)
			if err != nil {
				log.Error("error sending violation", "error", err)
				return err
			}
		}
	}
	log.Info("processing telemetry successfully")
	return nil

}

func hasDataChanged(previous, current *models.TelemetryData) bool {
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
		previous.RPM != current.RPM ||
		previous.Handbrake != current.Handbrake
}

func floatsEqual(a, b float64) bool {
	const epsilon = 0.000001
	return math.Abs(a-b) < epsilon
}

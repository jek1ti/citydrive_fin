package service

import (
	"context"
	"log/slog"

	"github.com/jekiti/citydrive/telemetry/internal/config"
	"github.com/jekiti/citydrive/telemetry/internal/models"
)

type ViolationService struct {
	config *config.ViolationsConfig
	log    *slog.Logger
}

func NewViolationService(cfg *config.ViolationsConfig, log *slog.Logger) *ViolationService {
	return &ViolationService{config: cfg, log: log}
}

func (s *ViolationService) CheckViolations(ctx context.Context, carID string, current *models.TelemetryData) []*models.Violation {
	traceID := ctx.Value("trace_id")
	log := s.log.With(
		"module", "violation.service",
		"function", "CheckViolations",
		"car_id", carID,
		"trace_id", traceID,
	)
	log.Info("checking violations")
	violations := []*models.Violation{}
	speedLimit := int32(s.config.SpeedLimit)
	fuelLimit := s.config.LowFuelLimit
	rpmLimit := int32(s.config.DriftRPMLimit)

	if current.Speed > speedLimit {
		violations = append(violations, createSpeedViolation(speedLimit, carID, current))
	}
	if current.RPM > rpmLimit && current.Handbrake {
		violations = append(violations, &models.Violation{
			Type:  models.ViolationTypeDrift,
			CarID: carID,
			Data:  *current,
		})
	}
	if current.Fuel < fuelLimit {
		violations = append(violations, &models.Violation{
			Type:  models.ViolationTypeLowFuel,
			CarID: carID,
			Data:  *current,
		})
	}
	if !current.Activated && !current.Locked && current.EngineOn && current.Speed != 0 {
		violations = append(violations, &models.Violation{
			Type:  models.ViolationTypeStealedAuto,
			CarID: carID,
			Data:  *current,
		})
	}
	log.Info("violation checked")
	return violations
}

func createSpeedViolation(speedLimit int32, carID string, current *models.TelemetryData) *models.Violation {
	if current.Speed-speedLimit < 20 {
		return &models.Violation{
			Type:  models.ViolationTypeSpeedingLow,
			CarID: carID,
			Data:  *current,
		}
	}
	if current.Speed-speedLimit > 20 && current.Speed-speedLimit <= 40 {
		return &models.Violation{
			Type:  models.ViolationTypeSpeedingMedium,
			CarID: carID,
			Data:  *current,
		}
	}
	return &models.Violation{
		Type:  models.ViolationTypeSpeedingHigh,
		CarID: carID,
		Data:  *current,
	}
}

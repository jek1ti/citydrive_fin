package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jekiti/citydrive/admin/internal/domain"
	"github.com/jekiti/citydrive/admin/internal/repository"
)

type Service interface {
	GetCarHistory(ctx context.Context, carID string, from, to int64) ([]domain.CarState, error)
	GetCarsHistory(ctx context.Context, from, to int64, activated *bool) (map[string][]domain.CarHistoryPoint, error)
	GetCarsNow(ctx context.Context) ([]domain.CarShort, error)
	GetCar(ctx context.Context, carID string) (domain.CarDetails, error)
}

type service struct {
	repoDB    repository.DBRepository
	repoCache repository.CacheRepository
	log       *slog.Logger
}

func NewService(repoDB repository.DBRepository, repoCache repository.CacheRepository, log *slog.Logger) Service {
	return &service{
		repoDB:    repoDB,
		repoCache: repoCache,
		log:       log,
	}
}

func (s *service) GetCarHistory(ctx context.Context, carID string, from, to int64) ([]domain.CarState, error) {
	log := s.log.With("module", "service", "function", "GetCarHistory", "car_id", carID)
	log.Info("fetching car history", "from", from, "to", to)
	if from >= to {
		log.Error("invalid time range: 'from' timestamp is greater than or equal to 'to' timestamp", "from", from, "to", to)
		return nil, domain.ErrInvalidTimeRange
	}
	states, err := s.repoDB.GetCarHistory(ctx, carID, from, to)
	if err != nil {
		log.Error("error fetching car history from repository", "error", err)
		return nil, err
	}

	var isStaying bool = true
	if len(states) == 0 || len(states) == 1 {
		log.Info("no car history found for the specified time range", "car_id", carID)
		return states, nil
	}

	for i := 1; i < len(states); i++ {
		if (states[i].Lon == states[i-1].Lon && states[i].Lat == states[i-1].Lat) || (states[i].Activated == false && states[i-1].Activated == false) {
			continue
		} else {
			isStaying = false
			break
		}
	}
	if isStaying {
		log.Info("car has not moved during the specified time range", "car_id", carID)
		return []domain.CarState{states[len(states)-1]}, nil
	}

	log.Info("successfully fetched car history", "count", len(states))
	return states, nil
}

func (s *service) GetCarsHistory(ctx context.Context, from, to int64, activated *bool) (map[string][]domain.CarHistoryPoint, error) {
	log := s.log.With("module", "service", "function", "GetCarsHistory")
	log.Info("fetching cars history", "from", from, "to", to, "activated", activated)
	if from >= to {
		log.Error("invalid time range: 'from' timestamp is greater than or equal to 'to' timestamp", "from", from, "to", to)
		return nil, domain.ErrInvalidTimeRange
	}
	if activated == nil {
		log.Error("activated is required, but is nil")
		return nil, fmt.Errorf("activated is required")
	}
	history, err := s.repoDB.GetCarsHistory(ctx, from, to, activated)
	if err != nil {
		log.Error("error fetching cars history from repository", "error", err)
		return nil, err
	}
	log.Info("successfully fetched cars history", "car_count", len(history))
	return history, nil
}

func (s *service) GetCarsNow(ctx context.Context) ([]domain.CarShort, error) {
	log := s.log.With("module", "service", "function", "GetActiveCars")
	log.Info("fetching active cars")
	cars, err := s.repoCache.GetCarsNow(ctx)
	if err != nil {
		log.Error("error fetching active cars from repository", "error", err)
	}
	log.Info("successfully fetched active cars", "count", len(cars))
	return cars, nil
}

func (s *service) GetCar(ctx context.Context, carID string) (domain.CarDetails, error) {
	log := s.log.With("module", "service", "function", "GetCarByID", "car_id", carID)
	log.Info("fetching active cars")
	if carID == "" {
		return domain.CarDetails{}, domain.ErrInvalidCarID
	}
	car, err := s.repoCache.GetCar(ctx, carID)
	if err != nil {
		log.Error("error fetching car by id from repository", "error", err)
		return domain.CarDetails{}, err
	}
	return car, nil
}

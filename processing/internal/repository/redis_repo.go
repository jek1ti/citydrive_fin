package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/jekiti/citydrive/processing/internal/config"
	"github.com/jekiti/citydrive/processing/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	SaveCarState(telemetry domain.CarTelemetry) error
	GetCarState(carID string) (*domain.CarTelemetry, error)
	Close() error
}

type RedisRepository struct {
	client *redis.Client
	log    *slog.Logger
	config *config.RedisConfig
}

func NewRedisRepository(cfg *config.RedisConfig, log *slog.Logger) (CacheRepository, error) {
	log = log.With("module", "repository", "function", "NewRedisRepository")
	log.Info("connecting to redis", "host", cfg.Host, "port", cfg.Port, "db", cfg.DB)
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Info("pinging redis to check connection")
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Error("failed to connect redis", "error", err)
		return nil, err
	}

	return &RedisRepository{
		client: client,
		log:    log,
		config: cfg,
	}, nil
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}

func (r *RedisRepository) SaveCarState(telemetry domain.CarTelemetry) error {
	log := r.log.With("function", "SaveCarState", "car_id", telemetry.CarID)

	if telemetry.CarID == "" {
        log.Warn("skip save: empty car_id")
        return nil
    }
	log.Info("marshaling car state to json")
	jsonData, err := json.Marshal(telemetry)
	if err != nil {
		log.Error("error marshaling car state to json", "error", err)
		return err
	}

	key := strings.Replace(r.config.KeyCarCurrent, "{car_id}", telemetry.CarID, 1)
	log.Info("saving car state to redis", "key", key)
	err = r.client.Set(context.Background(), key, jsonData, 0).Err()
	if err != nil {
		log.Error("error saving car state to redis", "error", err)
		return err
	}
	log.Info("car state saved to redis", "car_id", telemetry.CarID)
	return nil
}

func (r *RedisRepository) GetCarState(carID string) (*domain.CarTelemetry, error) {
	log := r.log.With("function", "GetCarState", "car_id", carID)

	key := strings.Replace(r.config.KeyCarCurrent, "{car_id}", carID, 1)
	data, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Info("no data found for carID", "car_id", carID)
			return nil, nil
		}
		log.Error("error getting car state from redis", "error", err)
		return nil, err
	}

	var telemetry domain.CarTelemetry
	err = json.Unmarshal([]byte(data), &telemetry)
	if err != nil {
		log.Error("error unmarshaling car state", "error", err)
		return nil, err
	}
	log.Info("car state retrieved from redis", "car_id", carID)
	return &telemetry, nil
}

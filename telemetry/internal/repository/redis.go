package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jekiti/citydrive/telemetry/internal/config"
	"github.com/jekiti/citydrive/telemetry/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	prefix string
	log    *slog.Logger
}

func NewRedisRepository(cfg *config.TelemetryConfig, log *slog.Logger) (*RedisRepository, error) {
	log = log.With("module", "repository", "function", "NewRedisRepository")
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		DialTimeout:  cfg.Redis.DialTimeout,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Error("failed to connect redis", "error", err)
		return nil, fmt.Errorf("redis conntection failed: %w", err)
	}
	return &RedisRepository{
		client: client,
		prefix: "car:state:",
		log: log,
	}, nil
}

func (r *RedisRepository) GetCarState(ctx context.Context, carID string) (*models.TelemetryData, error) {
	log := r.log.With("module", "repository", "function", "GetCarState")
	key := r.prefix + carID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Info("no data found for carID", "car_id", carID)
			return nil, nil
		}
		log.Error("error getting data from redis", "error", err)
		return nil, fmt.Errorf("error getting data from redis:%w", err)
	}

	var state models.TelemetryData
	json.Unmarshal([]byte(data), &state)
	return &state, nil
}

func (r *RedisRepository) SetCarState(ctx context.Context, carID string, data *models.TelemetryData) error {
	key := r.prefix + carID
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal car state: %w", err)
	}
	err = r.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set car state:%w", err)
	}
	return nil
}

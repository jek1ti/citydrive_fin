package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jekiti/citydrive/admin/internal/config"
	"github.com/jekiti/citydrive/admin/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	GetCarsNow(ctx context.Context) ([]domain.CarShort, error)
	GetCar(ctx context.Context, carID string) (domain.CarDetails, error)
	Close() error
}

type RedisRepository struct {
	client *redis.Client
	prefix string
	log    *slog.Logger
}

func NewRedisRepository(cfg *config.AdminConfig, log *slog.Logger) (CacheRepository, error) {
	log = log.With("module", "repository", "function", "NewRedisRepository")
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.URL,
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
		log:    log,
	}, nil
}

func (r *RedisRepository) GetCarsNow(ctx context.Context) ([]domain.CarShort, error) {
	log := r.log.With("module", "repository", "function", "GetActiveCars")
	var allKeys []string
	var cursor uint64 = 0

	for {
		log.Info("scanning keys from redis", "cursor", cursor)
		keys, nextCursor, err := r.client.Scan(ctx, cursor, r.prefix+"*", 100).Result()
		if err != nil {
			log.Error("error scanning keys from redis", "error", err)
			return nil, fmt.Errorf("error scanning keys from redis: %w", err)
		}
		allKeys = append(allKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	if len(allKeys) == 0 {
		log.Info("no active cars found")
		return nil, nil
	}

	log.Info("found active cars", "count", len(allKeys))
	data, err := r.client.MGet(ctx, allKeys...).Result()
	if err != nil {
		log.Error("error getting data from redis", "error", err)
		return nil, fmt.Errorf("error getting data from redis:%w", err)
	}

	var state []domain.CarShort

	for i, key := range allKeys {
		if data[i] == nil {
			continue
		}
		str, ok := data[i].(string)

		if !ok {
			log.Error("error asserting data to string", "key", key)
			continue
		}
		var car domain.CarDetails
		if err := json.Unmarshal([]byte(str), &car); err != nil {
			log.Error("error unmarshaling car data", "error", err, "key", key)
			continue
		}

		var carShort domain.CarShort
		if car.Activated{
			carID := strings.TrimPrefix(key, r.prefix)
			carShort = domain.CarShort{
				ID:    carID,
				Brand: car.Brand,
				Model: car.Model,
				Lat:   car.Lat,
				Lon:   car.Lon,
				Speed: car.Speed,
			}
			state = append(state, carShort)
		}
	}
	log.Info("active cars data retrieved from redis", "count", len(state))
	return state, nil
}

func (r *RedisRepository) GetCar(ctx context.Context, carID string) (domain.CarDetails, error) {
	log := r.log.With("module", "repository", "function", "GetCarByID", "car_id", carID)
	key := r.prefix + carID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Info("no data found for carID", "car_id", carID)
			return domain.CarDetails{}, domain.ErrCarNotFound
		}
		log.Error("error getting data from redis", "error", err)
		return domain.CarDetails{}, fmt.Errorf("error getting data from redis:%w", err)
	}
	var carDetails domain.CarDetails
	err = json.Unmarshal([]byte(data), &carDetails)
	if err != nil {
		log.Error("error unmarshaling car data", "error", err)
		return domain.CarDetails{}, fmt.Errorf("error unmarshaling car data: %w", err)
	}
	return carDetails, nil
}


func (r *RedisRepository) Close() error {
	return r.client.Close()
}
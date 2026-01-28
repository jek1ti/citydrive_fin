package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AdminConfig struct {
	GRPC     GRPCConfig
	Redis    RedisConfig
	Postgres PostgresConfig
	App      AppConfig
}

type GRPCConfig struct {
	Port string
}

type RedisConfig struct {
	URL          string
	Password     string
	DB           int
	PoolSize     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

type PostgresConfig struct {
	URL         string
	MaxConns    int
	ReadTimeout time.Duration
}

type AppConfig struct {
	Env      string
	LogLevel string
}

func LoadAdminConfig() *AdminConfig {
	_ = godotenv.Load()

	return &AdminConfig{
		GRPC: GRPCConfig{
			Port: getDefault("GRPC_PORT", "50053"),
		},
		Redis: RedisConfig{
			URL:          getDefault("REDIS_URL", "localhost:6379"),
			Password:     getDefault("REDIS_PASSWORD", ""),
			DB:           getIntDefault("REDIS_DB", 0),
			PoolSize:     getIntDefault("REDIS_POOL_SIZE", 10),
			ReadTimeout:  getDurationDefault("REDIS_READ_TIMEOUT", "3s"),
			WriteTimeout: getDurationDefault("REDIS_WRITE_TIMEOUT", "3s"),
			DialTimeout:  getDurationDefault("REDIS_DIAL_TIMEOUT", "3s"),
		},
		Postgres: PostgresConfig{
			URL:         mustGet("DB_URL"),
			MaxConns:    getIntDefault("DB_MAX_CONN", 10),
			ReadTimeout: getDurationDefault("DB_READ_TIMEOUT", "5s"),
		},
		App: AppConfig{
			Env:      getDefault("ENV", "development"),
			LogLevel: getDefault("LOG_LEVEL", "info"),
		},
	}
}

func (c *AdminConfig) Validate() error {
	if c.Redis.URL == "" {
		log.Fatal("REDIS_URL is required")
	}
	if c.Postgres.URL == "" {
		log.Fatal("DB_URL is required")
	}
	return nil
}

func getDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getIntDefault(key string, def int) int {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("error parsing int from env %s: %v", key, err)
	}
	return i
}

func getDurationDefault(key, def string) time.Duration {
	s := os.Getenv(key)
	if s == "" {
		d, _ := time.ParseDuration(def)
		return d
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("error parsing duration from env %s: %v", key, err)
	}
	return d
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(".env value is missing: %s", key)
	}
	return v
}

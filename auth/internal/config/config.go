package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	GRPCPort string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	MaxConn  int
	SSLMode  string
}

type JWTConfig struct {
	Algorithm  string
	SecretKey  string
	Expiration time.Duration
	PublicKeyPath string
	JWKSURL       string
}

type AppConfig struct {
	Env      string
	LogLevel string
}

func LoadConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Printf("Warning: could not load .env file from %s: %v", path, err)
	}
	return &Config{
		Server: ServerConfig{
			GRPCPort: mustGet("GRPC_PORT"),
		},
		Database: DatabaseConfig{
			URL:      mustGet("DB_URL"),
			Host:     mustGet("DB_HOST"),
			Port:     mustGetInt("DB_PORT"),
			Name:     mustGet("DB_NAME"),
			User:     mustGet("DB_USER"),
			Password: mustGet("DB_PASSWORD"),
			MaxConn:  mustGetInt("DB_MAX_CONN"),
			SSLMode:  mustGet("DB_SSL_MODE"),
		},
		JWT: JWTConfig{
			Algorithm:     getDefault("JWT_ALG", "HS256"),
			SecretKey:     mustGet("JWT_SECRET_KEY"),
			Expiration:    mustGetDuration("JWT_EXPIRATION"),
			PublicKeyPath: getDefault("JWT_PUBLIC_KEY_PATH", ""),
			JWKSURL:       getDefault("AUTH_JWKS_URL", ""),
		},
		App: AppConfig{
			Env:      mustGet("ENV"),
			LogLevel: mustGet("LOG_LEVEL"),
		},
	}
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(".env value is missing: %s", key)
	}
	return v
}

func getDefault(key, def string) string {
	v := os.Getenv(key)

	if v == "" {
		return def
	}

	return v
}
func mustGetInt(key string) int {
	s := mustGet(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("bad int %s: %v", key, err)
	}
	return i
}
func mustGetDuration(key string) time.Duration {
	s := mustGet(key)
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("bad duration %s: %v", key, err)
	}
	return d
}

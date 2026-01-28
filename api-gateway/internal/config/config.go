package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type GatewayConfig struct {
	HTTP    HTTPConfig
	GRPC    GRPCConfig
	JWT     JWTConfig
	App     AppConfig
	Tracing TracingConfig
}

type HTTPConfig struct {
	Port              string
	BodyLimitBytes    int
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type GRPCConfig struct {
	AuthAddr           string
	TelemetryAddr      string
	AdminAddr          string
	DialTimeout        time.Duration
	MaxCallRecvMsgSize int
}

type JWTConfig struct {
	SecretKey string
	CarSecretKey string
}

type AppConfig struct {
	Env      string
	LogLevel string
}

type TracingConfig struct {
	HeaderName string
}

func LoadGatewayConfig() *GatewayConfig {
	_ = godotenv.Load()
	return &GatewayConfig{
		HTTP: HTTPConfig{
			Port:              getDefault("HTTP_PORT", "8080"),
			BodyLimitBytes:    getIntDefault("GATEWAY_BODY_LIMIT_BYTES", 1048576),
			ReadHeaderTimeout: getDurationDefault("GATEWAY_READ_HEADER_TIMEOUT", "10s"),
			WriteTimeout:      getDurationDefault("GATEWAY_WRITE_TIMEOUT", "15s"),
			IdleTimeout:       getDurationDefault("GATEWAY_IDLE_TIMEOUT", "60s"),
		},
		GRPC: GRPCConfig{
			AuthAddr:           mustGet("AUTH_GRPC_ADDR"),
			TelemetryAddr:      mustGet("TELEMETRY_GRPC_ADDR"),
			AdminAddr:          mustGet("ADMIN_GRPC_ADDR"),
			DialTimeout:        getDurationDefault("GRPC_DIAL_TIMEOUT", "5s"),
			MaxCallRecvMsgSize: getIntDefault("GRPC_MAX_RECV_MSG_SIZE", 4194304),
		},
		JWT: JWTConfig{
			SecretKey: mustGet("JWT_SECRET_KEY"),
			CarSecretKey: mustGet("JWT_CAR_SECRET_KEY"),
		},
		App: AppConfig{
			Env:      getDefault("ENV", "development"),
			LogLevel: getDefault("LOG_LEVEL", "info"),
		},
		Tracing: TracingConfig{
			HeaderName: getDefault("TRACE_HEADER_NAME", "X-Trace-ID"),
		},
	}
}

func (c *GatewayConfig) Validate() error {
	if c.JWT.SecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is required")
	}
	if c.JWT.CarSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is required")
	}
	return nil
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
		log.Fatalf("error parsing int from env %s: %v", key, err)
	}
	return i
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

func mustGetDuration(key string) time.Duration {
	s := mustGet(key)
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("error parsing duration from env %s: %v", key, err)
	}
	return d
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

func getBoolDefault(key string, def bool) bool {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("error parsing bool from env %s: %v", key, err)
	}
	return b
}

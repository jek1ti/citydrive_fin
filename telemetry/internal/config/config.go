package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type TelemetryConfig struct {
	GRPC       GRPCConfig
	Redis      RedisConfig
	Kafka      KafkaConfig
	Violations ViolationsConfig
	App        AppConfig
	Processing ProcessingConfig
}

type GRPCConfig struct {
	Port                 string
	MaxConcurrentStreams uint32
	MaxRecvMsgSize       int
	ConnectionTimeout    time.Duration
}

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

type KafkaConfig struct {
	Brokers         []string
	ClientID        string
	TelemetryTopic  string
	ViolationsTopic string
	ProducerAcks    string
	Retries         int
	BatchSize       int
	BatchTimeout    time.Duration
	Compression     string
}

type ViolationsConfig struct {
	SpeedLimit    int
	SpeedMedium   int
	SpeedHigh     int
	DriftRPMLimit int
	LowFuelLimit  float64
}

type AppConfig struct {
	Env         string
	LogLevel    string
	ServiceName string
	MetricsPort string
}

type ProcessingConfig struct {
	WorkerPoolSize       int
	QueueSize            int
	StateChangeThreshold time.Duration
	MaxProcessingTime    time.Duration
}

func LoadTelemetryConfig() *TelemetryConfig {
	_ = godotenv.Load()

	telemetryTopic := os.Getenv("KAFKA_TOPIC_TELEMETRY_RAW")
	if telemetryTopic == "" {
		telemetryTopic = os.Getenv("KAFKA_TOPIC_TELEMETRY")
	}
	if telemetryTopic == "" {
		telemetryTopic = "telemetry.raw"
	}

	violationsTopic := os.Getenv("KAFKA_TOPIC_VIOLATIONS")
	if violationsTopic == "" {
		violationsTopic = os.Getenv("KAFKA_TOPIC_ALERTS")
	}
	if violationsTopic == "" {
		violationsTopic = "telemetry.violations"
	}

	return &TelemetryConfig{
		GRPC: GRPCConfig{
			Port:                 getDefault("GRPC_PORT", "50052"),
			MaxConcurrentStreams: uint32(getIntDefault("GRPC_MAX_CONCURRENT_STREAMS", 1000)),
			MaxRecvMsgSize:       getIntDefault("GRPC_MAX_RECV_MSG_SIZE", 10485760),
			ConnectionTimeout:    getDurationDefault("GRPC_CONNECTION_TIMEOUT", "10s"),
		},
		Redis: RedisConfig{
			Host:         getDefault("REDIS_HOST", "localhost"),
			Port:         getDefault("REDIS_PORT", "6379"),
			Password:     getDefault("REDIS_PASSWORD", ""),
			DB:           getIntDefault("REDIS_DB", 0),
			PoolSize:     getIntDefault("REDIS_POOL_SIZE", 100),
			ReadTimeout:  getDurationDefault("REDIS_READ_TIMEOUT", "3s"),
			WriteTimeout: getDurationDefault("REDIS_WRITE_TIMEOUT", "3s"),
			DialTimeout:  getDurationDefault("REDIS_DIAL_TIMEOUT", "5s"),
		},
		Kafka: KafkaConfig{
			Brokers:         getSliceDefault("KAFKA_BROKERS", []string{"localhost:9092"}),
			ClientID:        getDefault("KAFKA_CLIENT_ID", "telemetry-ingestion"),
			TelemetryTopic:  telemetryTopic,
			ViolationsTopic: violationsTopic,
			ProducerAcks:    getDefault("KAFKA_PRODUCER_ACKS", "all"),
			Retries:         getIntDefault("KAFKA_PRODUCER_RETRIES", 3),
			BatchSize:       getIntDefault("KAFKA_PRODUCER_BATCH_SIZE", 1000000),
			BatchTimeout:    getDurationDefault("KAFKA_PRODUCER_BATCH_TIMEOUT", "500ms"),
			Compression:     getDefault("KAFKA_PRODUCER_COMPRESSION", "lz4"),
		},
		Violations: ViolationsConfig{
			SpeedLimit:    getIntDefault("VIOLATION_SPEED_LIMIT", 110),
			SpeedMedium:   getIntDefault("VIOLATION_SPEED_MEDIUM", 130),
			SpeedHigh:     getIntDefault("VIOLATION_SPEED_HIGH", 150),
			DriftRPMLimit: getIntDefault("VIOLATION_DRIFT_RPM_LIMIT", 5000),
			LowFuelLimit:  getFloatDefault("VIOLATION_LOW_FUEL_LIMIT", 2.0),
		},
		App: AppConfig{
			Env:         getDefault("ENV", "development"),
			LogLevel:    getDefault("LOG_LEVEL", "info"),
			ServiceName: getDefault("SERVICE_NAME", "telemetry-ingestion"),
			MetricsPort: getDefault("METRICS_PORT", "9090"),
		},
		Processing: ProcessingConfig{
			WorkerPoolSize:       getIntDefault("TELEMETRY_WORKER_POOL_SIZE", 10),
			QueueSize:            getIntDefault("TELEMETRY_QUEUE_SIZE", 1000),
			StateChangeThreshold: getDurationDefault("STATE_CHANGE_THRESHOLD_MS", "5s"),
			MaxProcessingTime:    getDurationDefault("MAX_PROCESSING_TIME", "30s"),
		},
	}
}

func (c *TelemetryConfig) Validate() error {
	if len(c.Kafka.Brokers) == 0 {
		log.Fatal("KAFKA_BROKERS is required")
	}
	if c.Redis.Host == "" {
		log.Fatal("REDIS_HOST is required")
	}

	if c.Violations.SpeedLimit <= 0 {
		log.Fatal("VIOLATION_SPEED_LIMIT must be positive")
	}
	if c.Violations.LowFuelLimit < 0 || c.Violations.LowFuelLimit > 100 {
		log.Fatal("VIOLATION_LOW_FUEL_LIMIT must be between 0 and 100")
	}

	return nil
}

func getSliceDefault(key string, def []string) []string {
	s := os.Getenv(key)
	if s == "" {
		return def
	}

	var result []string
	start := 0
	for i, char := range s {
		if char == ',' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

func getFloatDefault(key string, def float64) float64 {
	s := os.Getenv(key)
	if s == "" {
		return def
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("error parsing float from env %s: %v", key, err)
	}
	return f
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

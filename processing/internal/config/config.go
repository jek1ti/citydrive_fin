package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ProcessorConfig struct {
	DB        DBConfig
	Redis     RedisConfig
	Kafka     KafkaConfig
	App       AppConfig
	Processor ProcessorSpecificConfig
}

type DBConfig struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	MaxConn  int
	SSLMode  string
}

type RedisConfig struct {
	Host             string
	Port             string
	Password         string
	DB               int
	PoolSize         int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	DialTimeout      time.Duration
	KeyCarCurrent    string
	KeyCarLastUpdate string
}

type KafkaConfig struct {
	Brokers         string
	TopicTelemetry  string
	ConsumerGroupID string
	AutoOffsetReset string
	ClientID        string
}

type AppConfig struct {
	Env         string
	LogLevel    string
	ServiceName string
	HTTPPort    string
}

type ProcessorSpecificConfig struct {
	BatchSize       int
	WorkerPoolSize  int
	CommitInterval  time.Duration
	PollTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func LoadProcessorConfig() *ProcessorConfig {
	_ = godotenv.Load()

	topicTelemetry := os.Getenv("KAFKA_TOPIC_TELEMETRY_RAW")
	if topicTelemetry == "" {
		topicTelemetry = os.Getenv("KAFKA_TOPIC_TELEMETRY")
	}
	if topicTelemetry == "" {
		log.Fatal(".env value is missing: KAFKA_TOPIC_TELEMETRY_RAW")
	}

	return &ProcessorConfig{
		DB: DBConfig{
			URL:      mustGet("DB_URL"),
			Host:     getDefault("DB_HOST", "localhost"),
			Port:     getDefault("DB_PORT", "5432"),
			Name:     mustGet("DB_NAME"),
			User:     mustGet("DB_USER"),
			Password: mustGet("DB_PASSWORD"),
			MaxConn:  getIntDefault("DB_MAX_CONN", 10),
			SSLMode:  getDefault("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:             getDefault("REDIS_HOST", "localhost"),
			Port:             getDefault("REDIS_PORT", "6379"),
			Password:         getDefault("REDIS_PASSWORD", ""),
			DB:               getIntDefault("REDIS_DB", 0),
			PoolSize:         getIntDefault("REDIS_POOL_SIZE", 100),
			ReadTimeout:      getDurationDefault("REDIS_READ_TIMEOUT", "3s"),
			WriteTimeout:     getDurationDefault("REDIS_WRITE_TIMEOUT", "3s"),
			DialTimeout:      getDurationDefault("REDIS_DIAL_TIMEOUT", "5s"),
			KeyCarCurrent:    getDefault("REDIS_KEY_CAR_CURRENT", "car:current:{car_id}"),
			KeyCarLastUpdate: getDefault("REDIS_KEY_CAR_LAST_UPDATE", "car:last_update:{car_id}"),
		},
		Kafka: KafkaConfig{
			Brokers:         mustGet("KAFKA_BROKERS"),
			TopicTelemetry:  topicTelemetry,
			ConsumerGroupID: getDefault("KAFKA_CONSUMER_GROUP_ID", "telemetry-processor-group"),
			AutoOffsetReset: getDefault("KAFKA_AUTO_OFFSET_RESET", "earliest"),
			ClientID:        getDefault("KAFKA_CLIENT_ID", "telemetry-processor"),
		},
		App: AppConfig{
			Env:         getDefault("ENV", "development"),
			LogLevel:    getDefault("LOG_LEVEL", "info"),
			ServiceName: getDefault("SERVICE_NAME", "telemetry-processor"),
			HTTPPort:    getDefault("HTTP_PORT", "8083"),
		},
		Processor: ProcessorSpecificConfig{
			BatchSize:       getIntDefault("PROCESSOR_BATCH_SIZE", 100),
			WorkerPoolSize:  getIntDefault("PROCESSOR_WORKER_POOL_SIZE", 5),
			CommitInterval:  getDurationDefault("PROCESSOR_COMMIT_INTERVAL", "5s"),
			PollTimeout:     getDurationDefault("PROCESSOR_POLL_TIMEOUT", "100ms"),
			ShutdownTimeout: getDurationDefault("PROCESSOR_SHUTDOWN_TIMEOUT", "30s"),
		},
	}
}

func (c *ProcessorConfig) Validate() error {
	if c.DB.URL == "" {
		log.Fatal("DB_URL is required")
	}
	if c.Kafka.Brokers == "" {
		log.Fatal("KAFKA_BROKERS is required")
	}
	if c.Kafka.TopicTelemetry == "" {
		log.Fatal("KAFKA_TOPIC_TELEMETRY is required")
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

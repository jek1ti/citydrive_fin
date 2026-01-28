package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jekiti/citydrive/telemetry/internal/config"
	"github.com/jekiti/citydrive/telemetry/internal/models"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	telemetryWriter  *kafka.Writer
	violationsWriter *kafka.Writer
	log              *slog.Logger
}

func NewKafkaProducer(cfg *config.TelemetryConfig, log *slog.Logger) (*KafkaProducer, error) {
	telemetryWriter, err := createWriter(cfg, "telemetry")
	if err != nil {
		return nil, fmt.Errorf("error creating telemetryKafkaWrite: %w", err)
	}
	violationsWriter, err := createWriter(cfg, "violations")
	if err != nil {
		return nil, fmt.Errorf("error creating telemetryKafkaWrite: %w", err)
	}

	return &KafkaProducer{
		telemetryWriter:  telemetryWriter,
		violationsWriter: violationsWriter,
		log:              log,
	}, nil
}

func createWriter(cfg *config.TelemetryConfig, name string) (*kafka.Writer, error) {

	brokers := cfg.Kafka.Brokers
	var topic string
	if name == "violations" {
		topic = cfg.Kafka.ViolationsTopic
	} else {
		topic = cfg.Kafka.TelemetryTopic
	}
	var requiredAcks kafka.RequiredAcks
	switch cfg.Kafka.ProducerAcks {
	case "all":
		requiredAcks = kafka.RequireAll
	case "one":
		requiredAcks = kafka.RequireOne
	default:
		requiredAcks = kafka.RequireAll
	}

	var compression kafka.Compression
	switch cfg.Kafka.Compression {
	case "lz4":
		compression = kafka.Lz4
	case "gzip":
		compression = kafka.Gzip
	default:
		compression = kafka.Lz4
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: requiredAcks,
		MaxAttempts:  cfg.Kafka.Retries,
		BatchSize:    cfg.Kafka.BatchSize,
		BatchTimeout: cfg.Kafka.BatchTimeout,
		Compression:  compression,
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		return nil, fmt.Errorf("kafka topic %s not available: %w", topic, err)
	}
	conn.Close()

	return writer, nil
}

func (p *KafkaProducer) SendTelemetry(ctx context.Context, carID string, data *models.TelemetryData) error {
	traceID := ctx.Value("trace_id")
	log := p.log.With(
		"module", "producer",
		"function", "SendTelemetry",
		"car_id", carID,
		"trace_id", traceID,
	)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("error sending telemetry", "error", err)
		return fmt.Errorf("failed sending telemetry:%w", err)
	}

	message := kafka.Message{
		Key:   []byte(carID),
		Value: jsonData,
	}
	err = p.telemetryWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Error("error writing message in telemetry topic", "error", err)
		return fmt.Errorf("failed to write message in telemetry Topic: %w", err)
	}
	log.Info("telemetry sended successfully")
	return nil
}

func (p *KafkaProducer) SendViolation(ctx context.Context, violation *models.Violation) error {
	traceID := ctx.Value("trace_id")
	log := p.log.With(
		"module", "producer",
		"function", "SendViolation",
		"car_id", violation.CarID,
		"trace_id", traceID,
	)
	jsonData, err := json.Marshal(violation)
	if err != nil {
		log.Error("error marshal violation", "error", err)
		return fmt.Errorf("failed marshal violation:%w", err)
	}

	message := kafka.Message{
		Key:   []byte(violation.Type),
		Value: jsonData,
	}
	err = p.violationsWriter.WriteMessages(ctx, message)
	if err != nil {
		log.Error("error writing message in violation Topic", "error", err)
		return fmt.Errorf("failed to write message in violation Topic: %w", err)
	}
	log.Info("violation sended successfully")
	return nil
}

func (p *KafkaProducer) Close() error {
	log := p.log.With(
		"module", "producer",
		"function", "Close",
	)
	err1 := p.telemetryWriter.Close()
	err2 := p.violationsWriter.Close()

	if err1 != nil || err2 != nil {
		log.Error("error closing producers", "error1", err1, "error2", err2)
		return fmt.Errorf("failed to close kafka producers: telemetry=%v, violations=%v", err1, err2)
	}
	log.Info("producer closed successfully")
	return nil
}

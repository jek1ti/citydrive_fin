package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jekiti/citydrive/processing/internal/config"
	"github.com/jekiti/citydrive/processing/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	GetMessages(ctx context.Context, count int) ([]domain.CarTelemetry, error)
	Commit() error
	Close() error
}

type KafkaConsumer struct {
	reader       *kafka.Reader
	config       *config.KafkaConfig
	lastMessages []kafka.Message
	log          *slog.Logger
}

func NewKafkaConsumer(config *config.KafkaConfig, log *slog.Logger) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(config.Brokers, ","),
		Topic:    config.TopicTelemetry,
		GroupID:  config.ConsumerGroupID,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  300 * time.Millisecond,
	})

	return &KafkaConsumer{
		reader:       reader,
		config:       config,
		log:          log,
		lastMessages: make([]kafka.Message, 0),
	}
}

func (kc *KafkaConsumer) GetMessages(ctx context.Context, count int) ([]domain.CarTelemetry, error) {
	log := kc.log.With("module", "repository", "function", "GetMessages")
	log.Info("reading messages from kafka", "count", count)
	messages := make([]domain.CarTelemetry, 0, count)

	newCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	for len(messages) < count {
		msg, err := kc.reader.ReadMessage(newCtx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				log.Info("read timeout", "collected", len(messages))
				break
			}
			log.Error("error reading message from kafka", "error", err)
			return nil, err
		}

		log.Info("read message from kafka", "offset", msg.Offset)

		var telemetry domain.CarTelemetry
		err = json.Unmarshal(msg.Value, &telemetry)
		if err != nil {
			log.Error("error unmarshaling message", "error", err)
		}
		telemetry.CarID = string(msg.Key)
		telemetry.ReceivedAt = msg.Time.Unix()

		messages = append(messages, telemetry)
		kc.lastMessages = append(kc.lastMessages, msg)
	}
	log.Info("getted messages successful", "msg", messages)
	return messages, nil
}

func (kc *KafkaConsumer) Commit() error {
	log := kc.log.With("module", "repository", "function", "Commit")
	log.Info("committing offsets to kafka", "count", len(kc.lastMessages))
	err := kc.reader.CommitMessages(context.Background(), kc.lastMessages...)
	if err != nil {
		log.Error("error committing messages", "error", err)
		return err
	}
	kc.lastMessages = nil
	return nil
}

func (kc *KafkaConsumer) Close() error {
	log := kc.log.With("module", "repository", "function", "Close")
	log.Info("closing kafka consumer")
	return kc.reader.Close()
}

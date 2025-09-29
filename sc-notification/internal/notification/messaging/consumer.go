package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"soul-connect/sc-notification/internal/config"
	"soul-connect/sc-notification/internal/notification/service"
)

type NotificationEvent struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

type NotificationConsumer struct {
	reader  *kafka.Reader
	service service.NotificationService
}

func NewNotificationConsumer(cfg *config.Config, svc service.NotificationService) (*NotificationConsumer, error) {
	if cfg.KafkaTopic == "" || cfg.KafkaBrokers == "" {
		return nil, errors.New("kafka configuration is incomplete")
	}

	brokers := parseBrokers(cfg.KafkaBrokers)
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are empty")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: cfg.KafkaGroupID,
		Topic:   cfg.KafkaTopic,
	})

	return &NotificationConsumer{
		reader:  reader,
		service: svc,
	}, nil
}

func (c *NotificationConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("notification consumer: fetch message error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if err := c.handleMessage(ctx, msg); err != nil {
			log.Printf("notification consumer: handle message error: %v", err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("notification consumer: commit message error: %v", err)
		}
	}
}

func (c *NotificationConsumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	var event NotificationEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return err
	}

	_, err = c.service.CreateNotification(ctx, userID, event.Content)
	return err
}

func (c *NotificationConsumer) Close() error {
	if c.reader == nil {
		return nil
	}
	return c.reader.Close()
}

func parseBrokers(raw string) []string {
	parts := strings.Split(raw, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}

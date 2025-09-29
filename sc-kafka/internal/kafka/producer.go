package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	segment "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"soul-connect/sc-kafka/internal/config"
)

// Producer is a thin wrapper around kafka-go writers for the configured topics.
type Producer struct {
	writers map[string]*segment.Writer
	topics  config.TopicConfig
}

// NewProducer constructs topic-specific writers based on the provided configuration.
func NewProducer(cfg config.Config) (*Producer, error) {
	if len(cfg.Broker.Addresses) == 0 {
		return nil, errors.New("at least one kafka broker address is required")
	}

	dialer := &segment.Dialer{
		Timeout:   cfg.Broker.DialTimeout,
		DualStack: true,
		ClientID:  cfg.Broker.ClientID,
	}

	if cfg.Broker.Username != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: cfg.Broker.Username,
			Password: cfg.Broker.Password,
		}
	}

	topics := []string{
		cfg.Topics.PostCreated,
		cfg.Topics.SubscriptionCreated,
		cfg.Topics.Notification,
	}

	writers := make(map[string]*segment.Writer, len(topics))
	for _, topic := range topics {
		if topic == "" {
			continue
		}

		writer := segment.NewWriter(segment.WriterConfig{
			Brokers:      cfg.Broker.Addresses,
			Topic:        topic,
			Dialer:       dialer,
			Balancer:     &segment.Hash{},
			BatchSize:    cfg.Producer.BatchSize,
			BatchTimeout: cfg.Producer.BatchTimeout,
			RequiredAcks: int(segment.RequireAll),
			ReadTimeout:  cfg.Broker.ReadTimeout,
			WriteTimeout: cfg.Broker.WriteTimeout,
		})

		writer.AllowAutoTopicCreation = cfg.Producer.AllowAutoTopicCreate
		writers[topic] = writer
	}

	if len(writers) == 0 {
		return nil, errors.New("no kafka writers were initialised")
	}

	return &Producer{writers: writers, topics: cfg.Topics}, nil
}

// Close closes all writers and aggregates potential errors.
func (p *Producer) Close() error {
	var err error
	for _, writer := range p.writers {
		if closeErr := writer.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	return err
}

// PublishPostCreated emits a post creation event to the configured topic.
func (p *Producer) PublishPostCreated(ctx context.Context, event PostCreatedEvent) error {
	key := keyOrFallback(event.ID, event.AuthorID)
	return p.publish(ctx, p.topics.PostCreated, key, event)
}

// PublishSubscriptionCreated emits a subscription creation event to the configured topic.
func (p *Producer) PublishSubscriptionCreated(ctx context.Context, event SubscriptionCreatedEvent) error {
	key := keyOrFallback(event.ID, event.SubscriberID)
	return p.publish(ctx, p.topics.SubscriptionCreated, key, event)
}

// PublishNotification emits a notification event to the configured topic.
func (p *Producer) PublishNotification(ctx context.Context, event NotificationEvent) error {
	key := keyOrFallback(event.ID, event.UserID)
	return p.publish(ctx, p.topics.Notification, key, event)
}

func (p *Producer) publish(ctx context.Context, topic string, key []byte, payload any) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("writer for topic %q is not configured", topic)
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	message := segment.Message{
		Key:   key,
		Value: value,
		Time:  time.Now().UTC(),
	}

	return writer.WriteMessages(ctx, message)
}

func keyOrFallback(primary, fallback string) []byte {
	if primary != "" {
		return []byte(primary)
	}

	if fallback != "" {
		return []byte(fallback)
	}

	return nil
}

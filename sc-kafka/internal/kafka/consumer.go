package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	segment "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"soul-connect/sc-kafka/internal/config"
)

// Consumer wires kafka-go readers for the configured topics and groups.
type Consumer struct {
	readers map[string]*segment.Reader
	topics  config.TopicConfig
}

// NewConsumer constructs a Consumer using the provided configuration.
func NewConsumer(cfg config.Config) (*Consumer, error) {
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

	topicsWithGroups := map[string]string{
		cfg.Topics.PostCreated:         cfg.Consumer.PostCreatedGroup,
		cfg.Topics.SubscriptionCreated: cfg.Consumer.SubscriptionGroup,
		cfg.Topics.Notification:        cfg.Consumer.NotificationGroup,
	}

	readers := make(map[string]*segment.Reader, len(topicsWithGroups))
	for topic, group := range topicsWithGroups {
		if topic == "" || group == "" {
			continue
		}

		readers[topic] = segment.NewReader(segment.ReaderConfig{
			Brokers:           cfg.Broker.Addresses,
			GroupID:           group,
			Topic:             topic,
			Dialer:            dialer,
			MinBytes:          cfg.Consumer.MinBytes,
			MaxBytes:          cfg.Consumer.MaxBytes,
			CommitInterval:    cfg.Consumer.CommitInterval,
			HeartbeatInterval: cfg.Consumer.HeartbeatInterval,
			SessionTimeout:    cfg.Consumer.SessionTimeout,
			ReadLagInterval:   -1,
		})
	}

	if len(readers) == 0 {
		return nil, errors.New("no kafka readers were initialised")
	}

	return &Consumer{readers: readers, topics: cfg.Topics}, nil
}

// Close closes all kafka readers and aggregates potential errors.
func (c *Consumer) Close() error {
	var err error
	for _, reader := range c.readers {
		if closeErr := reader.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	return err
}

// ConsumePostCreated runs the provided handler for every post-created event.
func (c *Consumer) ConsumePostCreated(ctx context.Context, handler func(context.Context, PostCreatedEvent) error) error {
	return c.consumePostCreated(ctx, handler)
}

// ConsumeSubscriptionCreated runs the handler for subscription-created events.
func (c *Consumer) ConsumeSubscriptionCreated(ctx context.Context, handler func(context.Context, SubscriptionCreatedEvent) error) error {
	return c.consumeSubscriptionCreated(ctx, handler)
}

// ConsumeNotification runs the handler for notification events.
func (c *Consumer) ConsumeNotification(ctx context.Context, handler func(context.Context, NotificationEvent) error) error {
	return c.consumeNotification(ctx, handler)
}

func (c *Consumer) consumePostCreated(ctx context.Context, handler func(context.Context, PostCreatedEvent) error) error {
	return consumeMessages(ctx, c.readers, c.topics.PostCreated, handler)
}

func (c *Consumer) consumeSubscriptionCreated(ctx context.Context, handler func(context.Context, SubscriptionCreatedEvent) error) error {
	return consumeMessages(ctx, c.readers, c.topics.SubscriptionCreated, handler)
}

func (c *Consumer) consumeNotification(ctx context.Context, handler func(context.Context, NotificationEvent) error) error {
	return consumeMessages(ctx, c.readers, c.topics.Notification, handler)
}

// DrainTopic discards messages until the provided context is cancelled.
// It can be useful for manual topic resets during development.
func (c *Consumer) DrainTopic(ctx context.Context, topic string) error {
	reader, ok := c.readers[topic]
	if !ok {
		return fmt.Errorf("reader for topic %q is not configured", topic)
	}

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("fetch message: %w", err)
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka: failed to commit drained message on topic %s: %v", topic, err)
		}
	}
}

func consumeMessages[T any](ctx context.Context, readers map[string]*segment.Reader, topic string, handler func(context.Context, T) error) error {
	reader, ok := readers[topic]
	if !ok {
		return fmt.Errorf("reader for topic %q is not configured", topic)
	}

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("fetch message: %w", err)
		}

		var event T
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka: failed to decode message on topic %s: %v", topic, err)
			if commitErr := reader.CommitMessages(context.Background(), msg); commitErr != nil {
				log.Printf("kafka: failed to commit corrupted message on topic %s: %v", topic, commitErr)
			}
			continue
		}

		if err := handler(ctx, event); err != nil {
			log.Printf("kafka: handler returned error for topic %s: %v", topic, err)
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka: failed to commit message on topic %s: %v", topic, err)
		}
	}
}

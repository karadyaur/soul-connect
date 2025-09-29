package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, key string, payload []byte) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

type noopProducer struct{}

func NewProducer(brokers []string, topic string) (Producer, error) {
	if len(brokers) == 0 || topic == "" {
		return &noopProducer{}, nil
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	return &kafkaProducer{writer: writer}, nil
}

func (p *kafkaProducer) Publish(ctx context.Context, key string, payload []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: payload})
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}

func (p *noopProducer) Publish(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (p *noopProducer) Close() error {
	return nil
}

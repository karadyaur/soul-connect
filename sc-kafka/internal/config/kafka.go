package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config describes the kafka broker and the producer/consumer clients
// that are required by the service.
type Config struct {
	Broker   BrokerConfig
	Topics   TopicConfig
	Producer ProducerClientConfig
	Consumer ConsumerClientConfig
}

// BrokerConfig contains low level connection options for Kafka brokers.
type BrokerConfig struct {
	Addresses    []string
	ClientID     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Username     string
	Password     string
}

// TopicConfig holds the topic names that correspond to the
// domain events processed by the service.
type TopicConfig struct {
	PostCreated         string
	SubscriptionCreated string
	Notification        string
}

// ProducerClientConfig exposes producer-specific settings.
type ProducerClientConfig struct {
	BatchSize            int
	BatchTimeout         time.Duration
	AllowAutoTopicCreate bool
}

// ConsumerClientConfig exposes consumer-specific settings including
// the consumer group identifiers for each domain event.
type ConsumerClientConfig struct {
	PostCreatedGroup  string
	SubscriptionGroup string
	NotificationGroup string
	MinBytes          int
	MaxBytes          int
	CommitInterval    time.Duration
	HeartbeatInterval time.Duration
	SessionTimeout    time.Duration
}

// Load initialises the configuration using a .env file located at path
// (if present) and environment variables. Missing optional values fall back
// to sensible defaults.
func Load(path string) (Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(filepath.Join(path, ".env"))
		if err := v.ReadInConfig(); err != nil {
			var configFileErr viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileErr) {
				return Config{}, fmt.Errorf("read config: %w", err)
			}
		}
	}

	v.AutomaticEnv()

	setDefaults(v)

	cfg := Config{
		Broker: BrokerConfig{
			Addresses:    parseBrokers(v.GetString("KAFKA_BROKERS")),
			ClientID:     v.GetString("KAFKA_CLIENT_ID"),
			DialTimeout:  v.GetDuration("KAFKA_DIAL_TIMEOUT"),
			ReadTimeout:  v.GetDuration("KAFKA_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("KAFKA_WRITE_TIMEOUT"),
			Username:     v.GetString("KAFKA_USERNAME"),
			Password:     v.GetString("KAFKA_PASSWORD"),
		},
		Topics: TopicConfig{
			PostCreated:         v.GetString("KAFKA_TOPIC_POST_CREATED"),
			SubscriptionCreated: v.GetString("KAFKA_TOPIC_SUBSCRIPTION_CREATED"),
			Notification:        v.GetString("KAFKA_TOPIC_NOTIFICATION"),
		},
		Producer: ProducerClientConfig{
			BatchSize:            v.GetInt("KAFKA_PRODUCER_BATCH_SIZE"),
			BatchTimeout:         v.GetDuration("KAFKA_PRODUCER_BATCH_TIMEOUT"),
			AllowAutoTopicCreate: v.GetBool("KAFKA_PRODUCER_AUTO_CREATE_TOPIC"),
		},
		Consumer: ConsumerClientConfig{
			PostCreatedGroup:  v.GetString("KAFKA_CONSUMER_POST_CREATED_GROUP"),
			SubscriptionGroup: v.GetString("KAFKA_CONSUMER_SUBSCRIPTION_GROUP"),
			NotificationGroup: v.GetString("KAFKA_CONSUMER_NOTIFICATION_GROUP"),
			MinBytes:          v.GetInt("KAFKA_CONSUMER_MIN_BYTES"),
			MaxBytes:          v.GetInt("KAFKA_CONSUMER_MAX_BYTES"),
			CommitInterval:    v.GetDuration("KAFKA_CONSUMER_COMMIT_INTERVAL"),
			HeartbeatInterval: v.GetDuration("KAFKA_CONSUMER_HEARTBEAT_INTERVAL"),
			SessionTimeout:    v.GetDuration("KAFKA_CONSUMER_SESSION_TIMEOUT"),
		},
	}

	if len(cfg.Broker.Addresses) == 0 {
		return Config{}, errors.New("kafka brokers are not configured")
	}

	if cfg.Topics.PostCreated == "" || cfg.Topics.SubscriptionCreated == "" || cfg.Topics.Notification == "" {
		return Config{}, errors.New("kafka topics are not fully configured")
	}

	if cfg.Producer.BatchSize <= 0 {
		cfg.Producer.BatchSize = 1
	}

	if cfg.Consumer.MinBytes <= 0 {
		cfg.Consumer.MinBytes = 1
	}

	if cfg.Consumer.MaxBytes <= 0 {
		cfg.Consumer.MaxBytes = 10 << 20 // 10 MiB
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("KAFKA_BROKERS", "kafka-broker:9092")
	v.SetDefault("KAFKA_CLIENT_ID", "sc-kafka")
	v.SetDefault("KAFKA_DIAL_TIMEOUT", "5s")
	v.SetDefault("KAFKA_READ_TIMEOUT", "10s")
	v.SetDefault("KAFKA_WRITE_TIMEOUT", "10s")

	v.SetDefault("KAFKA_TOPIC_POST_CREATED", "post.created")
	v.SetDefault("KAFKA_TOPIC_SUBSCRIPTION_CREATED", "subscription.created")
	v.SetDefault("KAFKA_TOPIC_NOTIFICATION", "notification.created")

	v.SetDefault("KAFKA_PRODUCER_BATCH_SIZE", 100)
	v.SetDefault("KAFKA_PRODUCER_BATCH_TIMEOUT", "200ms")
	v.SetDefault("KAFKA_PRODUCER_AUTO_CREATE_TOPIC", true)

	v.SetDefault("KAFKA_CONSUMER_POST_CREATED_GROUP", "sc-post-consumers")
	v.SetDefault("KAFKA_CONSUMER_SUBSCRIPTION_GROUP", "sc-subscription-consumers")
	v.SetDefault("KAFKA_CONSUMER_NOTIFICATION_GROUP", "sc-notification-consumers")
	v.SetDefault("KAFKA_CONSUMER_MIN_BYTES", 1)
	v.SetDefault("KAFKA_CONSUMER_MAX_BYTES", 10<<20)
	v.SetDefault("KAFKA_CONSUMER_COMMIT_INTERVAL", "1s")
	v.SetDefault("KAFKA_CONSUMER_HEARTBEAT_INTERVAL", "3s")
	v.SetDefault("KAFKA_CONSUMER_SESSION_TIMEOUT", "30s")
}

func parseBrokers(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}

	return out
}

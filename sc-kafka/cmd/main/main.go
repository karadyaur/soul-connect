package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"soul-connect/sc-kafka/internal/config"
	appkafka "soul-connect/sc-kafka/internal/kafka"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("kafka: failed to load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	consumer, err := appkafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("kafka: failed to initialise consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("kafka: failed to close consumers cleanly: %v", err)
		}
	}()

	log.Println("kafka: starting domain event workers")

	var wg sync.WaitGroup
	runWorker := func(name string, fn func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("kafka: worker %s stopped with error: %v", name, err)
			}
		}()
	}

	runWorker("post-created", func(ctx context.Context) error {
		return consumer.ConsumePostCreated(ctx, func(ctx context.Context, event appkafka.PostCreatedEvent) error {
			log.Printf("kafka: received post created event: %+v", event)
			return nil
		})
	})

	runWorker("subscription-created", func(ctx context.Context) error {
		return consumer.ConsumeSubscriptionCreated(ctx, func(ctx context.Context, event appkafka.SubscriptionCreatedEvent) error {
			log.Printf("kafka: received subscription created event: %+v", event)
			return nil
		})
	})

	runWorker("notification", func(ctx context.Context) error {
		return consumer.ConsumeNotification(ctx, func(ctx context.Context, event appkafka.NotificationEvent) error {
			log.Printf("kafka: received notification event: %+v", event)
			return nil
		})
	})

	<-ctx.Done()
	log.Println("kafka: shutting down workers")
	wg.Wait()
}

func loadConfig() (config.Config, error) {
	configPath := os.Getenv("KAFKA_CONFIG_PATH")
	return config.Load(configPath)
}

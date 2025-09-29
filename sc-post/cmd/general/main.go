package main

import (
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"soul-connect/sc-post/internal/config"
	"soul-connect/sc-post/internal/events"
	"soul-connect/sc-post/internal/server"
	"soul-connect/sc-post/internal/services"
	"soul-connect/sc-post/pkg/database"
	"soul-connect/sc-post/pkg/kafka"
	postpb "soul-connect/sc-post/pkg/postpb"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	pool, err := database.NewPostgresDB(&cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer pool.Close()

	brokers := parseBrokers(cfg.KafkaBrokers)
	producer, err := kafka.NewProducer(brokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer producer.Close()

	publisher := events.NewPostEventPublisher(producer, cfg.KafkaTopic)
	svc := services.NewServices(pool, publisher)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	postpb.RegisterPostServiceServer(grpcServer, server.NewPostServer(svc))

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.ServerPort, err)
	}

	log.Printf("post service gRPC server listening on %s", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func parseBrokers(value string) []string {
	if value == "" {
		return nil
	}
	items := strings.Split(value, ",")
	brokers := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}

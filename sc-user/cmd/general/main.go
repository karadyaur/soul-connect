package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"soul-connect/sc-user/internal/config"
	"soul-connect/sc-user/internal/generated"
	"soul-connect/sc-user/internal/server"
	"soul-connect/sc-user/internal/services"
	"soul-connect/sc-user/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	pool, err := database.NewPostgresDB(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port %s: %v", cfg.ServerPort, err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	service := services.NewService(pool)
	userServer := server.NewUserServer(service)

	generated.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("User gRPC server listening on port %s", cfg.ServerPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

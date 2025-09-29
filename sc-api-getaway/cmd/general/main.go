package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"soul-connect/sc-api-getaway/internal/config"
	"soul-connect/sc-api-getaway/internal/controllers"
	"soul-connect/sc-api-getaway/internal/generated"
	"soul-connect/sc-api-getaway/internal/routers"
	"syscall"
	"time"
)

func main() {
	newConfig, err := config.LoadConfig("./")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config: %v\n", err)
		os.Exit(1)
	}

	conn, err := grpc.Dial(newConfig.GrpcAuthTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to Auth Service: %v", err)
	}
	defer conn.Close()

	authServiceClient := generated.NewAuthServiceClient(conn)
	newController := controllers.NewController(authServiceClient)

	newRouter := routers.NewRouter(&newConfig, newController)
	newRouter.SetRoutes()

	newServer := &http.Server{
		Addr:    ":" + newConfig.ServerPort,
		Handler: newRouter.Gin,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Server is running on port %s\n", newConfig.ServerPort)
		if err := newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start failed: %s\n", err)
		}
	}()

	// Set up signal catching
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown initiated...")

	// Context for graceful shutdown with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := newServer.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Error: %v", err)
	}

	// Waiting for the shutdown context to be done or timeout
	<-ctx.Done()
	log.Println("Server shutdown completed or timed out")

	log.Println("Server exiting")
	os.Exit(0)
}

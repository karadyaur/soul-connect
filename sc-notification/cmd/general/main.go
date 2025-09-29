package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"soul-connect/sc-notification/internal/config"
	db "soul-connect/sc-notification/internal/db/sqlc"
	httpapi "soul-connect/sc-notification/internal/notification/api/http"
	"soul-connect/sc-notification/internal/notification/messaging"
	"soul-connect/sc-notification/internal/notification/repository"
	"soul-connect/sc-notification/internal/notification/service"
	"soul-connect/sc-notification/pkg/database"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	pool, err := database.NewPostgresDB(&cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	repo := repository.NewNotificationRepository(queries)
	svc := service.NewNotificationService(repo)

	router := gin.Default()
	handler := httpapi.NewNotificationHandler(svc)
	handler.RegisterRoutes(router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var consumer *messaging.NotificationConsumer
	if cfg.KafkaBrokers != "" && cfg.KafkaTopic != "" {
		consumer, err = messaging.NewNotificationConsumer(&cfg, svc)
		if err != nil {
			log.Printf("notification kafka consumer disabled: %v", err)
		} else {
			go consumer.Start(ctx)
			defer func() {
				if err := consumer.Close(); err != nil {
					log.Printf("notification kafka consumer close error: %v", err)
				}
			}()
		}
	} else {
		log.Println("notification kafka consumer not configured")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("notification server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("notification server shutdown error: %v", err)
	}
}

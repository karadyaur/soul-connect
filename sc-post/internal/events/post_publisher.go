package events

import (
	"context"
	"encoding/json"
	"time"

	"soul-connect/sc-post/internal/models"
	"soul-connect/sc-post/pkg/kafka"
)

type PostCreatedEvent struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type PostEventPublisher interface {
	PublishPostCreated(ctx context.Context, post models.Post) error
}

type postEventPublisher struct {
	producer kafka.Producer
	topic    string
}

func NewPostEventPublisher(producer kafka.Producer, topic string) PostEventPublisher {
	if producer == nil {
		return &noopPostPublisher{}
	}
	return &postEventPublisher{producer: producer, topic: topic}
}

func (p *postEventPublisher) PublishPostCreated(ctx context.Context, post models.Post) error {
	if p.topic == "" {
		return nil
	}
	event := PostCreatedEvent{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		CreatedAt: post.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.producer.Publish(ctx, post.ID, payload)
}

type noopPostPublisher struct{}

func (n *noopPostPublisher) PublishPostCreated(context.Context, models.Post) error {
	return nil
}

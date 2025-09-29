package kafka

import "time"

// PostCreatedEvent describes an event emitted when a new post is created.
type PostCreatedEvent struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// SubscriptionCreatedEvent represents a new subscription between two users.
type SubscriptionCreatedEvent struct {
	ID           string    `json:"id"`
	SubscriberID string    `json:"subscriber_id"`
	CreatorID    string    `json:"creator_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// NotificationEvent is a lightweight message that should be delivered to a user.
type NotificationEvent struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Type      string            `json:"type"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

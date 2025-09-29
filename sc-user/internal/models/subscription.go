package models

type ModifySubscriptionInput struct {
	SubscriberID string `json:"subscriber_id"`
	AuthorID     string `json:"author_id"`
}

type SubscriptionList struct {
	SubscriberID string   `json:"subscriber_id"`
	AuthorIDs    []string `json:"author_ids"`
}

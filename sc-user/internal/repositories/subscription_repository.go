package repositories

import (
	"context"

	db "soul-connect/sc-user/internal/db/sqlc"
)

type ISubscriptionRepository interface {
	Subscribe(ctx context.Context, subscriberID, authorID string) error
	Unsubscribe(ctx context.Context, subscriberID, authorID string) error
	ListAuthorIDs(ctx context.Context, subscriberID string) ([]string, error)
}

type SubscriptionRepository struct {
	queries db.Querier
}

func NewSubscriptionRepository(queries db.Querier) *SubscriptionRepository {
	return &SubscriptionRepository{queries: queries}
}

func (r *SubscriptionRepository) Subscribe(ctx context.Context, subscriberID, authorID string) error {
	subscriber, err := stringToUUID(subscriberID)
	if err != nil {
		return err
	}
	author, err := stringToUUID(authorID)
	if err != nil {
		return err
	}

	_, err = r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		SubscriberID: subscriber,
		AuthorID:     author,
	})
	return err
}

func (r *SubscriptionRepository) Unsubscribe(ctx context.Context, subscriberID, authorID string) error {
	subscriber, err := stringToUUID(subscriberID)
	if err != nil {
		return err
	}
	author, err := stringToUUID(authorID)
	if err != nil {
		return err
	}

	return r.queries.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
		SubscriberID: subscriber,
		AuthorID:     author,
	})
}

func (r *SubscriptionRepository) ListAuthorIDs(ctx context.Context, subscriberID string) ([]string, error) {
	subscriber, err := stringToUUID(subscriberID)
	if err != nil {
		return nil, err
	}

	authors, err := r.queries.GetSubscriptionsByUserID(ctx, subscriber)
	if err != nil {
		return nil, err
	}

	return uuidSliceToStrings(authors), nil
}

var _ ISubscriptionRepository = (*SubscriptionRepository)(nil)

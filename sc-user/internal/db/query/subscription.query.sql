-- name: CreateSubscription :one
INSERT INTO subscriptions (subscriber_id, author_id)
VALUES (@subscriber_id, @author_id)
RETURNING id, subscriber_id, author_id, created_at;

-- name: GetSubscriptionsByUserID :many
SELECT author_id
FROM subscriptions
WHERE subscriber_id = @subscriber_id;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE subscriber_id = @subscriber_id AND author_id = @author_id;
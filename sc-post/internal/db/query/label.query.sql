-- name: AddLabelToPost :exec
INSERT INTO labels_posts (label_id, post_id)
VALUES (@label_id, @post_id)
RETURNING label_id, post_id;

-- name: GetLabelsForPost :many
SELECT l.id, l.name
FROM labels_posts lp
    JOIN labels l ON lp.label_id = l.id
WHERE lp.post_id = @post_id;

-- name: RemoveLabelFromPost :exec
DELETE FROM labels_posts
WHERE label_id = @label_id AND post_id = @post_id;

-- name: GetAllLabels :many
SELECT id, name
FROM labels;

-- name: GetPostsByLabel :many
SELECT p.id, p.user_id, p.title, p.description, p.likes_count, p.created_at, p.updated_at
FROM posts p
    JOIN labels_posts lp ON p.id = lp.post_id
WHERE lp.label_id = @label_id;
-- name: CreatePost :one
INSERT INTO posts (user_id, title, description)
VALUES (@user_id, @title, @description)
RETURNING id, user_id, title, description, likes_count, created_at, updated_at;

-- name: GetPostByID :one
SELECT id, user_id, title, description, likes_count, created_at, updated_at
FROM posts
WHERE id = @id
LIMIT 1;

-- name: UpdatePost :exec
UPDATE posts
SET title = COALESCE(@title, title),
    description = COALESCE(@description, description),
    updated_at = NOW()
WHERE id = @id;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = @id;

-- name: GetPostsWithCommentsAndLikes :many
SELECT
    p.id AS post_id,
    p.user_id AS post_user_id,
    p.title AS post_title,
    p.description AS post_description,
    p.likes_count AS post_likes,
    COUNT(DISTINCT c.id) AS total_comments,
    COUNT(DISTINCT l.id) AS total_likes
FROM posts p
         LEFT JOIN comments c ON p.id = c.post_id
         LEFT JOIN likes l ON p.id = l.post_id
GROUP BY p.id
ORDER BY p.created_at DESC;
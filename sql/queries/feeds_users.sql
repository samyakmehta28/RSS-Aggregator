-- name: CreateFeedFollowUser :one
INSERT INTO feeds_users (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowUser :many
SELECT * FROM feeds_users WHERE user_id = $1;

-- name: DeleteFeedFollowUser :exec
DELETE FROM feeds_users WHERE user_id = $1 AND feed_id = $2;
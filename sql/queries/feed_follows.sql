-- name: FollowFeed :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE id = $1;

-- name: GetFollowedFeeds :many
SELECT * FROM feed_follows
WHERE user_id = $1;
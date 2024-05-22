// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const followFeed = `-- name: FollowFeed :one
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, feed_id, created_at, updated_at
`

type FollowFeedParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, followFeed,
		arg.ID,
		arg.UserID,
		arg.FeedID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FeedID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFollowedFeeds = `-- name: GetFollowedFeeds :many
SELECT id, user_id, feed_id, created_at, updated_at FROM feed_follows
WHERE user_id = $1
`

func (q *Queries) GetFollowedFeeds(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowedFeeds, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FeedID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE id = $1
`

func (q *Queries) UnfollowFeed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, id)
	return err
}

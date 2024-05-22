-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT user_fk
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT feed_fk
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_togethet_con
    UNIQUE(id, user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
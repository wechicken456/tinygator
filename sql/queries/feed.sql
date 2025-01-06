-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE feeds.name=$1;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = current_timestamp, updated_at = current_timestamp WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;

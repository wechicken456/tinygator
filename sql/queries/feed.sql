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



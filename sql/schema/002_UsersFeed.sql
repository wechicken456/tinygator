-- +goose Up
SET timezone = 'UTC';
CREATE TABLE users
	(
		id			UUID		NOT NULL,
		created_at 	TIMESTAMP NOT NULL,
		updated_at 	TIMESTAMP NOT NULL,
		name		VARCHAR(30) NOT NULL,
		CONSTRAINT USERS_PK
			PRIMARY KEY(id)
	);

CREATE TABLE feeds (
    id			UUID		NOT NULL,
    created_at 	TIMESTAMP NOT NULL,
    updated_at 	TIMESTAMP NOT NULL,
    name		VARCHAR(30) NOT NULL,
    url            TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT FEEDS_PK
        PRIMARY KEY(id),
    CONSTRAINT FEEDS_USERS_FK
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows
    (
        id			UUID		NOT NULL,
        created_at 	TIMESTAMP NOT NULL,
        updated_at 	TIMESTAMP NOT NULL,
        user_id		UUID		NOT NULL,
        feed_id		UUID		NOT NULL,
        CONSTRAINT FEED_FOLLOWS_PK
            PRIMARY KEY(id),
        CONSTRAINT FEED_FOLLOWS_UNIQUE
            UNIQUE(user_id, feed_id),
        CONSTRAINT FEED_FOLLOWS_USERS_FK
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
        CONSTRAINT FEED_FOLLOWS_FEEDS_FK
            FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
    );


-- +goose Down
DROP TABLE users;
DROP TABLE feeds;
DROP TABLE feed_follows;
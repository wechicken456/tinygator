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

-- +goose Down
DROP TABLE users;
DROP TABLE feeds;
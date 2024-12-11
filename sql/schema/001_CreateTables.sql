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

-- +goose Down
DROP TABLE users;

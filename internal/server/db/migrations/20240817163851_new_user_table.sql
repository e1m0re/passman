-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(255)       NOT NULL,
    password VARCHAR(255)       NOT NULL,
    UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS users_data
(
    id       SERIAL PRIMARY KEY NOT NULL,
    "user"   INT                NOT NULL REFERENCES users,
    metadata varchar(255)       NOT NULL,
    file     varchar(255)       NOT NULL,
    checksum varchar(32)        NOT NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE users;
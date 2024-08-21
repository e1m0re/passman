-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users
(
    id       SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(255)       NOT NULL,
    password VARCHAR(255)       NOT NULL,
    UNIQUE (username)
);

CREATE TABLE users_data_items
(
    id       SERIAL PRIMARY KEY NOT NULL,
    "type"   INT                NOT NULL,
    "user"   INT                NOT NULL REFERENCES users,
    file     varchar(255)       NOT NULL,
    checksum varchar(32)        NOT NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE users;
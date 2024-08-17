-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY                           NOT NULL,
    name       VARCHAR(255)                                 NOT NULL,
    password   VARCHAR(255)                                 NOT NULL,
    email      VARCHAR(255)                                 NOT NULL,
    created_ad TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE users;
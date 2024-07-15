-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table roles;
-- +goose StatementEnd

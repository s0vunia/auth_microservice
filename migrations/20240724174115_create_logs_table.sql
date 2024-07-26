-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs
(
    id         SERIAL PRIMARY KEY,
    message    VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd

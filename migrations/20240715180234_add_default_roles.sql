-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name) VALUES ('user') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('admin') ON CONFLICT (name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name IN ('user', 'admin');
-- +goose StatementEnd

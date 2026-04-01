-- +goose Up
-- +goose StatementBegin
ALTER TABLE members ADD COLUMN qbo_name VARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE members DROP COLUMN qbo_name;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_email_idx;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE INDEX members_name_idx ON members (name);
CREATE INDEX member_aliases_member_id_name_idx ON member_aliases (member_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX member_aliases_member_id_name_idx;
DROP INDEX members_name_idx;
-- +goose StatementEnd

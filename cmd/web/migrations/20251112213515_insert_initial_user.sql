-- +goose Up
-- +goose StatementBegin
INSERT INTO users (first_name, last_name, email, password, access_level, created_at, updated_at) 
VALUES ('admin', 'user', 'admin@admin.com', '$2a$12$WGoSeh49OcJuOwcrKhFh7OrbZ7xu9EKwGehKJLwg.vtAUvn8VrxCC', 1, '2025-10-01', '2025-10-01');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users 
WHERE email = 'admin@admin.com';
-- +goose StatementEnd

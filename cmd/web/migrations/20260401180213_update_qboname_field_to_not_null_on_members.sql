-- +goose Up
-- +goose StatementBegin

-- Step 2: Populate the new column with data for existing rows
-- Replace 'default_value' or a calculated value as needed
UPDATE members SET qbo_name = '' WHERE qbo_class IS NULL;

-- Step 3: Alter the column to be NOT NULL (required)
ALTER TABLE members ALTER COLUMN qbo_name SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE members ALTER COLUMN qbo_name DROP NOT NULL;
-- +goose StatementEnd

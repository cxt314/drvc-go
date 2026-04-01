-- +goose Up
-- +goose StatementBegin
-- Step 1: Add a new, nullable column
ALTER TABLE vehicles ADD COLUMN qbo_class VARCHAR(255) NULL;

-- Step 2: Populate the new column with data for existing rows
-- Replace 'default_value' or a calculated value as needed
UPDATE vehicles SET qbo_class = name WHERE qbo_class IS NULL;

-- Step 3: Alter the column to be NOT NULL (required)
ALTER TABLE vehicles ALTER COLUMN qbo_class SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vehicles DROP COLUMN qbo_class;
-- +goose StatementEnd

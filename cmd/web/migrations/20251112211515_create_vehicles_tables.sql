-- +goose Up
-- +goose StatementBegin
CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    make VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    fuel_type VARCHAR(2) NOT NULL,
    purchase_price INTEGER DEFAULT 0,
    purchase_date DATE, 
    vin VARCHAR(255) DEFAULT '',
    license_plate VARCHAR(255) DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    sale_price INTEGER DEFAULT 0,
    sale_date DATE,
    billing_type VARCHAR(255) NOT NULL,
    base_per_mile INTEGER DEFAULT 0,
    secondary_per_mile INTEGER DEFAULT 0,
    minimum_fee INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vehicles;
-- +goose StatementEnd

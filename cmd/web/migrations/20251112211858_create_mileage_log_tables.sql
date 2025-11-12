-- +goose Up
-- +goose StatementBegin
CREATE TABLE mileage_logs (
    id SERIAL PRIMARY KEY,
    vehicle_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    start_odometer INTEGER NOT NULL,
    end_odometer INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles (id) ON DELETE CASCADE
);

CREATE TABLE trips (
    id SERIAL PRIMARY KEY,
    mileage_log_id INTEGER NOT NULL,
    trip_date DATE,
    start_mileage INTEGER NOT NULL,
    end_mileage INTEGER NOT NULL,
    destination VARCHAR(255) NOT NULL,
    purpose VARCHAR(255) NOT NULL,
    long_distance_days INTEGER DEFAULT 0 NOT NULL,
    billing_rate VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (mileage_log_id) REFERENCES mileage_logs (id) ON DELETE CASCADE
);

CREATE TABLE riders (
    id SERIAL PRIMARY KEY,
    trip_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (trip_id) REFERENCES trips (id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE riders;
DROP TABLE trips;
DROP TABLE mileage_logs;
-- +goose StatementEnd

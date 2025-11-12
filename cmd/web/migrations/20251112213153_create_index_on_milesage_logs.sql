-- +goose Up
-- +goose StatementBegin
CREATE INDEX mileage_logs_vehicle_id_idx ON mileage_logs (vehicle_id);
CREATE INDEX mileage_logs_vehicle_year_month_idx ON mileage_logs (vehicle_id, year, month);
CREATE INDEX trips_mileage_log_id_idx ON trips (mileage_log_id);
CREATE INDEX riders_trip_member_idx ON riders (trip_id, member_id);
CREATE INDEX mileage_logs_year_month_idx ON mileage_logs (year, month);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX mileage_logs_vehicle_year_month_idx;
DROP INDEX mileage_logs_vehicle_id_idx;
DROP INDEX trips_mileage_log_id_idx;
DROP INDEX riders_trip_member_idx;
DROP INDEX mileage_logs_year_month_idx;
-- +goose StatementEnd

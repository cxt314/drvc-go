package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// vehicleCols lists the columns in the vehicles table EXCEPT "id"
const vehicleCols = `name, year, make, model, fuel_type,
				purchase_price, purchase_date, vin, license_plate,
				is_active, sale_price, sale_date,
				billing_type, base_per_mile, secondary_per_mile, minimum_fee,
				created_at, updated_at`

// InsertVehicle inserts a Vehicle into the database
func (m *postgresDBRepo) InsertVehicle(v models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	stmt := fmt.Sprintf(`INSERT INTO vehicles (%s)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		vehicleCols)

	_, err := m.DB.ExecContext(ctx, stmt,
		v.Name, v.Year, v.Make, v.Model, v.FuelType,
		v.PurchasePrice, v.PurchaseDate, v.Vin, v.LicensePlate,
		v.Active, v.SalePrice, v.SaleDate,
		v.BillingType, v.BasePerMile, v.SecondaryPerMile, v.MinimumFee,
		time.Now(), time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// scanRowsToVehicles takes a pointer to *sql.Rows and scans those values into a slice of Vehicles
func scanRowsToVehicles(rows *sql.Rows) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	for rows.Next() {
		m := models.Vehicle{}
		err := rows.Scan(&m.ID, &m.Name, &m.Year, &m.Make, &m.Model, &m.FuelType,
			&m.PurchasePrice, &m.PurchaseDate, &m.Vin, &m.LicensePlate,
			&m.Active, &m.SalePrice, &m.SaleDate,
			&m.BillingType, &m.BasePerMile, &m.SecondaryPerMile, &m.MinimumFee,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return vehicles, err
		}

		vehicles = append(vehicles, m)
	}
	err := rows.Err()
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

// AllVehicles returns a slice of all vehicles in database
func (m *postgresDBRepo) AllVehicles() ([]models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM vehicles ORDER BY name`, vehicleCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of vehicles to return
	return scanRowsToVehicles(rows)
}

// GetVehicleByActive returns a slice of all vehicles that have status = active
func (m *postgresDBRepo) GetVehicleByActive(active bool) ([]models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM vehicles WHERE is_active=$1 ORDER BY name`, vehicleCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of vehicles to return
	return scanRowsToVehicles(rows)
}

// GetVehicleByID returns one vehicle from a given id
func (m *postgresDBRepo) GetVehicleByID(id int) (models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.Vehicle

	q := fmt.Sprintf(`SELECT id, %s FROM vehicles
		WHERE id = $1`, vehicleCols)

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan single db row into vehicle model
	err := row.Scan(&v.ID, &v.Name, &v.Year, &v.Make, &v.Model, &v.FuelType,
		&v.PurchasePrice, &v.PurchaseDate, &v.Vin, &v.LicensePlate,
		&v.Active, &v.SalePrice, &v.SaleDate,
		&v.BillingType, &v.BasePerMile, &v.SecondaryPerMile, &v.MinimumFee,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return v, err
	}

	return v, nil
}

// UpdateVehicle updates a vehicle in the database
func (m *postgresDBRepo) UpdateVehicle(v models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE vehicles SET
			name = $1,
			year = $2,
			make = $3,
			model = $4,
			fuel_type = $5,
			purchase_price = $6,
			purchase_date = $7,
			vin = $8,
			license_plate = $9,
			is_active = $10,
			sale_price = $11,
			sale_date = $12,
			billing_type = $13,
			base_per_mile = $14,
			secondary_per_mile = $15,
			minimum_fee = $16,
			updated_at = $17
		WHERE id =  $18
		`

	_, err := m.DB.ExecContext(ctx, q,
		v.Name,
		v.Year,
		v.Make,
		v.Model,
		v.FuelType,
		v.PurchasePrice,
		v.PurchaseDate,
		v.Vin,
		v.LicensePlate,
		v.Active,
		v.SalePrice,
		v.SaleDate,
		v.BillingType,
		v.BasePerMile,
		v.SecondaryPerMile,
		v.MinimumFee,
		time.Now(),
		v.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateVehicleActiveByID updates the active status of a vehicle by id
func (m *postgresDBRepo) UpdateVehicleActiveByID(id int, active bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE vehicles SET
			is_active = $1,
			updated_at = $2
		WHERE id =  $3
		`

	_, err := m.DB.ExecContext(ctx, q,
		active,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteVehicle deletes one vehicle by id
// We should almost never actually delete a vehicle
// because it will be referenced by a lot of mileage logs.
func (m *postgresDBRepo) DeleteVehicle(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `DELETE from vehicles WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

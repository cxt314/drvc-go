package dbrepo

import (
	"context"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertVehicle inserts a Vehicle into the database
func (m *postgresDBRepo) InsertVehicle(v models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO vehicles (name, year, make, model, fuel_type,
				purchase_price, purchase_date, vin, license_plate,
				created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := m.DB.ExecContext(ctx, stmt,
		v.Name,
		v.Year,
		v.Make,
		v.Model,
		v.FuelType,
		v.PurchasePrice,
		v.PurchaseDate,
		v.Vin,
		v.LicensePlate,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) AllVehicles() ([]models.Vehicle, error) {
	q := `SELECT id,
				name,
				year,
				make,
				model,
				fuel_type,
				purchase_price,
				purchase_date,
				vin,
				license_plate,
				created_at,
				updated_at
		FROM vehicles`

	// execute our DB query
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of vehicles to return
	var vehicles []models.Vehicle

	for rows.Next() {
		m := models.Vehicle{}
		err = rows.Scan(&m.ID, &m.Name, &m.Year, &m.Make, &m.Model, &m.FuelType,
			&m.PurchasePrice, &m.PurchaseDate, &m.Vin, &m.LicensePlate,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return vehicles, err
		}

		vehicles = append(vehicles, m)
	}
	err = rows.Err()
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

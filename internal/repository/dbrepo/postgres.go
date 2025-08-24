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

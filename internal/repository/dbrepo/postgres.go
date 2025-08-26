package dbrepo

import (
	"context"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

const contextTimeout = 3 * time.Second

// AllUsers returns a slice of all users in the database
func (m *postgresDBRepo) AllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `SELECT first_name, last_name, email, password, access_level, created_at, updated_at
		FROM users`

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of users to return
	var users []models.User

	for rows.Next() {
		m := models.User{}
		err = rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Email, &m.Password, &m.AccessLevel,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return users, err
		}

		users = append(users, m)
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}

	return users, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(v models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE users SET
			first_name = $1,
			last_name = $2,
			email = $3,
			access_level = $4,
			updated_at = $5
		`

	_, err := m.DB.ExecContext(ctx, q,
		v.FirstName,
		v.LastName,
		v.Email,
		v.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// InsertVehicle inserts a Vehicle into the database
func (m *postgresDBRepo) InsertVehicle(v models.Vehicle) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
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

// AllVehicles returns a slice of all vehicles in database
func (m *postgresDBRepo) AllVehicles() ([]models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

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
	rows, err := m.DB.QueryContext(ctx, q)
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

// GetVehicleByID returns one vehicle from a given id
func (m *postgresDBRepo) GetVehicleByID(id int) (models.Vehicle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.Vehicle

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
		FROM vehicles
		WHERE id = $1`

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan db row into vehicle model
	err := row.Scan(&v.ID, &v.Name, &v.Year, &v.Make, &v.Model, &v.FuelType,
		&v.PurchasePrice, &v.PurchaseDate, &v.Vin, &v.LicensePlate,
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
			updated_at = $10
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
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteVehicle deletes one vehicle by id
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

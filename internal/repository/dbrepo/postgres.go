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

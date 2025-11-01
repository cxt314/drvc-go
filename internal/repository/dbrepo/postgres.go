package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const contextTimeout = 3 * time.Second

// AllUsers returns a slice of all users in the database
func (m *postgresDBRepo) AllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
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

// GetUserByID returns a user by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
		FROM users WHERE id=$1`

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan results into user
	u := models.User{}

	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.AccessLevel,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
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
		WHERE id = $6
		`

	_, err := m.DB.ExecContext(ctx, q,
		v.FirstName,
		v.LastName,
		v.Email,
		v.AccessLevel,
		time.Now(),
		v.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// generatePasswordHash takes a password string and returns a hashed password string
func generatePasswordHash(pw string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pw), 12)

	return string(hashedPassword)
}

// InsertUser inserts a User into the database
func (m *postgresDBRepo) InsertUser(v models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	stmt := `INSERT INTO users (first_name, last_name, email, password, access_level, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

	hashedPassword := generatePasswordHash(v.Password)

	_, err := m.DB.ExecContext(ctx, stmt,
		v.FirstName, v.LastName, v.Email, hashedPassword, v.AccessLevel,
		time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

// UpdateUserPassword updates the password for a given user
func (m *postgresDBRepo) UpdateUserPassword(v models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE users SET
			password = $1,
			updated_at = $2
		WHERE id = $3
		`
	hashedPassword := generatePasswordHash(v.Password)

	_, err := m.DB.ExecContext(ctx, q,
		hashedPassword,
		time.Now(),
		v.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	// successfully authenticated
	return id, hashedPassword, nil
}

// runInTx takes a func and runs it in a transaction
func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

// runInTxReturnID takes a func and runs it in a transaction, returning the ID of the newly inserted item
func runInTxReturnID(db *sql.DB, fn func(tx *sql.Tx) (int, error)) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	insertedID, err := fn(tx)
	if err == nil {
		return insertedID, tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return 0, errors.Join(err, rollbackErr)
	}

	return 0, err
}

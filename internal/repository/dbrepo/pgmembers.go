package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// memberCols lists the columns in the members table EXCEPT "id"
const memberCols = `name, email, created_at, updated_at`
const aliasCols = `member_id, name, created_at, updated_at`

// InsertMember inserts a Member into the database
func (m *postgresDBRepo) InsertMember(v models.Member) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		var lastInsertId int
		// insert into members table & return inserted member id
		stmt := fmt.Sprintf(`INSERT INTO members (%s)
				VALUES ($1, $2, $3, $4)
				RETURNING id`,
			memberCols)

		err := tx.QueryRowContext(ctx, stmt,
			v.Name, v.Email,
			time.Now(), time.Now(),
		).Scan(&lastInsertId)
		if err != nil {
			return err
		}

		// insert aliases into member_aliases table
		for _, a := range v.Aliases {
			stmt := fmt.Sprintf(`INSERT INTO member_aliases (%s)
				VALUES ($1, $2, $3, $4)`,
				aliasCols)

			_, err := tx.ExecContext(ctx, stmt,
				lastInsertId, a.Name,
				time.Now(), time.Now(),
			)

			if err != nil {
				return err
			}

		}

		return nil
	})
}

// scanRowsToMembers takes a pointer to *sql.Rows and scans those values into a slice of Members
func scanRowsToMembers(rows *sql.Rows) ([]models.Member, error) {
	var members []models.Member

	for rows.Next() {
		m := models.Member{}
		err := rows.Scan(&m.ID, &m.Name, &m.Email,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return members, err
		}

		members = append(members, m)
	}
	err := rows.Err()
	if err != nil {
		return members, err
	}

	return members, nil
}

// AllMembers returns a slice of all members in database
func (m *postgresDBRepo) AllMembers() ([]models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM members`, memberCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of members to return
	return scanRowsToMembers(rows)
}

// GetMemberByActive returns a slice of all members that have status = active
func (m *postgresDBRepo) GetMemberByActive(active bool) ([]models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM members WHERE is_active=$1`, memberCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of members to return
	return scanRowsToMembers(rows)
}

// GetMemberByID returns one member from a given id
func (m *postgresDBRepo) GetMemberByID(id int) (models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.Member

	q := fmt.Sprintf(`SELECT id, %s FROM members
		WHERE id = $1`, memberCols)

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan single db row into member model
	err := row.Scan(&v.ID, &v.Name, &v.Email,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return v, err
	}

	return v, nil
}

// UpdateMember updates a member in the database
func (m *postgresDBRepo) UpdateMember(v models.Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE members SET
			name = $1,
			email = $2,
			updated_at = $3
		WHERE id =  $4
		`

	_, err := m.DB.ExecContext(ctx, q,
		v.Name,
		v.Email,
		time.Now(),
		v.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMemberActiveByID updates the active status of a member by id
func (m *postgresDBRepo) UpdateMemberActiveByID(id int, active bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `UPDATE members SET
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

// DeleteMember deletes one member by id
// We should almost never actually delete a member
// because it will be referenced by a lot of mileage logs.
func (m *postgresDBRepo) DeleteMember(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := `DELETE from members WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

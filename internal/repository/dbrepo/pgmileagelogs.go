package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// memberCols lists the columns in the members table EXCEPT "id"
const mileageLogCols = `created_at, updated_at`
const TripCols = `created_at, updated_at`

// InsertMileageLog inserts a MileageLog into the database. This is wrapped in a transaction
// due to needing to insert a member and possible member aliaes
func (m *postgresDBRepo) InsertMileageLog(v models.MileageLog) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		var lastInsertId int
		// insert into members table & return inserted member id
		stmt := fmt.Sprintf(`INSERT INTO members (%s)
				VALUES ($1, $2, $3, $4)
				RETURNING id`,
			mileageLogCols)

		err := tx.QueryRowContext(ctx, stmt,
			v.Name, 
			time.Now(), time.Now(),
		).Scan(&lastInsertId)
		if err != nil {
			return err
		}

		// insert aliases into member_aliases table
		for _, a := range v.Trips {
			err := insertTripsTx(tx, ctx, lastInsertId, a.Name)

			if err != nil {
				return err
			}

		}

		return nil
	})
}

func insertTripsTx(tx *sql.Tx, ctx context.Context, member_id int, alias string) error {
	stmt := fmt.Sprintf(`INSERT INTO member_aliases (%s)
				VALUES ($1, $2, $3, $4)`,
		aliasCols)

	_, err := tx.ExecContext(ctx, stmt,
		member_id, alias,
		time.Now(), time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// scanRowsToMileageLogs takes a pointer to *sql.Rows and scans those values into a slice of MileageLogs
func scanRowsToMileageLogs(rows *sql.Rows) ([]models.MileageLog, error) {
	var members []models.MileageLog

	for rows.Next() {
		m := models.MileageLog{}
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

// scanRowsToMileageLogsTrips takes a pointer to *sql.Rows and scans those values into a slice of Tripes
func scanRowsToTripes(rows *sql.Rows) ([]models.Trip, error) {
	var aliases []models.Trip
	var memberId int

	for rows.Next() {
		m := models.Trip{}
		err := rows.Scan(&m.ID, &memberId, &m.Name,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return aliases, err
		}

		aliases = append(aliases, m)
	}
	err := rows.Err()
	if err != nil {
		return aliases, err
	}

	return aliases, nil
}

// AllMileageLogs returns a slice of all members in database. Does not populate member aliases
func (m *postgresDBRepo) AllMileageLogs() ([]models.MileageLog, error) {
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
	return scanRowsToMileageLogs(rows)
}

// GetMileageLogByActive returns a slice of all members that have status = active. Does not populate member aliases
func (m *postgresDBRepo) GetMileageLogByActive(active bool) ([]models.MileageLog, error) {
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
	return scanRowsToMileageLogs(rows)
}

// GetMileageLogByID returns one member from a given id, populates aliases
func (m *postgresDBRepo) GetMileageLogByID(id int) (models.MileageLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.MileageLog

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

	// get member aliases
	q = fmt.Sprintf(`SELECT id, %s FROM member_aliases WHERE member_id = $1`, aliasCols)
	rows, err := m.DB.QueryContext(ctx, q, id)
	if err != nil {
		return v, err
	}
	defer rows.Close()

	v.Trips, err = scanRowsToTripes(rows)
	if err != nil {
		return v, err
	}

	return v, nil
}

// UpdateMileageLog updates a member in the database
func (m *postgresDBRepo) UpdateMileageLog(v models.MileageLog) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		q := `UPDATE members SET
			name = $1,
			email = $2,
			updated_at = $3
		WHERE id =  $4 `

		_, err := tx.ExecContext(ctx, q,
			v.Name,
			v.Email,
			time.Now(),
			v.ID,
		)
		if err != nil {
			return err
		}

		// delete and re-add member_aliases. This avoids needing to check for updated aliases
		// delete existing member_aliases
		q = `DELETE from member_aliases WHERE member_id = $1`
		_, err = tx.ExecContext(ctx, q, v.ID)
		if err != nil {
			return err
		}

		// re-add member_aliases
		for _, a := range v.Trips {
			err := insertTripsTx(tx, ctx, v.ID, a.Name)

			if err != nil {
				return err
			}

		}

		return nil
	})
}

// UpdateMileageLogActiveByID updates the active status of a member by id
func (m *postgresDBRepo) UpdateMileageLogActiveByID(id int, active bool) error {
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

// DeleteMileageLog deletes one member by id. Also deletes all member_aliases with that member id
// We should almost never actually delete a member
// because it will be referenced by a lot of mileage logs.
func (m *postgresDBRepo) DeleteMileageLog(id int) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		q := `DELETE from member_aliases WHERE member_id = $1`
		_, err := tx.ExecContext(ctx, q, id)
		if err != nil {
			return err
		}

		q = `DELETE from members WHERE id=$1`

		_, err = tx.ExecContext(ctx, q, id)
		if err != nil {
			return err
		}

		return nil
	})
}

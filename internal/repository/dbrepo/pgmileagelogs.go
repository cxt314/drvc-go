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

// UpdateMileageLog updates a mielage log in the database
func (m *postgresDBRepo) UpdateMileageLog(v models.MileageLog) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		q := `UPDATE mileage_logs SET
			vehicle_id = $1,
			name = $2,
			year = $3,
			month = $4, 
			start_odometer = $5,
			end_odometer = $6,
			updated_at = $7
		WHERE id =  $8 `

		_, err := tx.ExecContext(ctx, q,
			v.Vehicle.ID,
			v.Name,
			v.Year,
			v.Month,
			v.StartOdometer,
			v.EndOdometer,
			time.Now(),
			v.ID,
		)
		if err != nil {
			return err
		}

		// do we process trips here?

		return nil
	})
}

// DeleteMileageLog deletes one mielage log by id. Also deletes all trips with that mileage log id
func (m *postgresDBRepo) DeleteMileageLog(id int) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		q := `DELETE from trips WHERE mileage_log_id = $1`
		_, err := tx.ExecContext(ctx, q, id)
		if err != nil {
			return err
		}

		q = `DELETE from mileage_logs WHERE id=$1`

		_, err = tx.ExecContext(ctx, q, id)
		if err != nil {
			return err
		}

		return nil
	})
}

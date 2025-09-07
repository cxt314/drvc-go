package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// memberCols lists the columns in the members table EXCEPT "id"
const mileageLogCols = `vehicle_id, name, year, month, start_odometer, end_odometer, created_at, updated_at`
const tripCols = `mileage_log_id, trip_date, start_mileage, end_mileage, is_long_distance, destination,
		purpose, created_at, updated_at`
const riderCols = `trip_id, member_id, created_at, updated_at`

// InsertMileageLog inserts a MileageLog into the database. This is wrapped in a transaction
// due to needing to insert trips and riders as well
func (m *postgresDBRepo) InsertMileageLog(v models.MileageLog) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		var lastInsertId int
		// insert into members table & return inserted member id
		stmt := fmt.Sprintf(`INSERT INTO mileage_logs (%s)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id`,
			mileageLogCols)

		err := tx.QueryRowContext(ctx, stmt,
			v.Vehicle.ID, v.Name, v.Year, v.Month, v.StartOdometer, v.EndOdometer,
			time.Now(), time.Now(),
		).Scan(&lastInsertId)
		if err != nil {
			return err
		}

		// insert trips into trips table
		for _, a := range v.Trips {
			err := insertTripsTx(tx, ctx, lastInsertId, a)

			if err != nil {
				return err
			}

		}

		return nil
	})
}

// insertTripsTx takes a transaction and inserts a Trip into the database
func insertTripsTx(tx *sql.Tx, ctx context.Context, mileage_log_id int, v models.Trip) error {
	stmt := fmt.Sprintf(`INSERT INTO trips (%s)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		aliasCols)

	_, err := tx.ExecContext(ctx, stmt,
		mileage_log_id, v.TripDate, v.StartMileage, v.EndMileage,
		v.IsLongDistance, v.Destination, v.Purpose,
		time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}

	// TODO: insert riders for a trip

	return nil
}

// scanRowsToMileageLogs takes a pointer to *sql.Rows and scans those values into a slice of MileageLogs
func scanRowsToMileageLogs(rows *sql.Rows) ([]models.MileageLog, error) {
	var logs []models.MileageLog

	for rows.Next() {
		m := models.MileageLog{}
		err := rows.Scan(&m.ID, &m.Vehicle.ID, &m.Name, &m.Year, &m.Month,
			&m.StartOdometer, &m.EndOdometer,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return logs, err
		}

		logs = append(logs, m)
	}
	err := rows.Err()
	if err != nil {
		return logs, err
	}

	return logs, nil
}

// scanRowsToTrips takes a pointer to *sql.Rows and scans those values into a slice of Trips
// also gets the slice of riders for the trip
func (m *postgresDBRepo) scanRowsToTrips(rows *sql.Rows) ([]models.Trip, error) {
	var trips []models.Trip
	var mileage_log_id int

	for rows.Next() {
		t := models.Trip{}
		err := rows.Scan(&t.ID, &mileage_log_id, &t.TripDate, &t.StartMileage,
			&t.EndMileage, &t.IsLongDistance, &t.Destination, &t.Purpose,
			&t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return trips, err
		}

		riders, err := m.getRidersByTripID(t.ID)
		if err != nil {
			return trips, err
		}

		t.Riders = riders
		trips = append(trips, t)
	}
	err := rows.Err()
	if err != nil {
		return trips, err
	}

	return trips, nil
}

// AllMileageLogs returns a slice of all mileage logs in database. Does not populate trips
func (m *postgresDBRepo) AllMileageLogs() ([]models.MileageLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM mileage_logs`, mileageLogCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of members to return
	return scanRowsToMileageLogs(rows)
}

// GetMileageLogByVehicleID returns a slice of all members that have status = active. Does not populate member aliases
func (m *postgresDBRepo) GetMileageLogsByVehicleID(vehicle_id int) ([]models.MileageLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM mileage_logs WHERE vehicle_id=$1`, mileageLogCols)

	// execute our DB query
	rows, err := m.DB.QueryContext(ctx, q, vehicle_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process query result into slice of members to return
	return scanRowsToMileageLogs(rows)
}

// GetMileageLogByID returns one mileage_log from a given id, populates trips
func (m *postgresDBRepo) GetMileageLogByID(id int) (models.MileageLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.MileageLog

	q := fmt.Sprintf(`SELECT id, %s FROM mileage_logs
		WHERE id = $1`, mileageLogCols)

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan single db row into mileage log model
	err := row.Scan(&v.ID, &v.Vehicle.ID, &v.Name, &v.Year,
		&v.Month, &v.StartOdometer, &v.EndOdometer,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return v, err
	}

	// get trips
	q = fmt.Sprintf(`SELECT id, %s FROM trips WHERE mileage_log_id = $1`, tripCols)
	rows, err := m.DB.QueryContext(ctx, q, id)
	if err != nil {
		return v, err
	}
	defer rows.Close()

	v.Trips, err = m.scanRowsToTrips(rows)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (m *postgresDBRepo) getRidersByTripID(id int) ([]models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var riders []models.Member

	q := `SELECT member_id FROM riders WHERE trip_id = $1`
	rows, err := m.DB.QueryContext(ctx, q, id)
	if err != nil {
		return riders, err
	}
	defer rows.Close()

	for rows.Next() {
		var member_id int

		err := rows.Scan(&member_id)
		if err != nil {
			return riders, err
		}

		// get member by id
		member, err := m.GetMemberByID(member_id)
		if err != nil {
			return riders, err
		}

		riders = append(riders, member)
	}
	err = rows.Err()
	if err != nil {
		return riders, err
	}

	return riders, nil
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

		// delete riders first
		// 		select *
		// FROM mileage_logs m
		// JOIN trips t on t.mileage_log_id = m.id
		// JOIN riders r on r.trip_id = t.id
		// WHERE m.id = 1

		// select *
		// from riders r
		// join trips t on r.trip_id = t.id
		// where t.mileage_log_id = 1

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

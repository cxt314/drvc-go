package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// memberCols lists the columns in the members table EXCEPT "id"
const memberCols = `name, email, is_active, created_at, updated_at`
const aliasCols = `member_id, name, created_at, updated_at`

// InsertMember inserts a Member into the database. This is wrapped in a transaction
// due to needing to insert a member and possible member aliaes
func (m *postgresDBRepo) InsertMember(v models.Member) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		var lastInsertId int
		// insert into members table & return inserted member id
		stmt := fmt.Sprintf(`INSERT INTO members (%s)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`,
			memberCols)

		err := tx.QueryRowContext(ctx, stmt,
			v.Name, v.Email, v.Active,
			time.Now(), time.Now(),
		).Scan(&lastInsertId)
		if err != nil {
			return err
		}

		// insert aliases into member_aliases table
		for _, a := range v.Aliases {
			err := insertMemberAliasesTx(tx, ctx, lastInsertId, a.Name)
			/*stmt := fmt.Sprintf(`INSERT INTO member_aliases (%s)
				VALUES ($1, $2, $3, $4)`,
				aliasCols)

			_, err := tx.ExecContext(ctx, stmt,
				lastInsertId, a.Name,
				time.Now(), time.Now(),
			)
			*/
			if err != nil {
				return err
			}

		}

		return nil
	})
}

// insertMemberAliasesTx is a helper function that takes a transaction and uses it to insert member aliases
func insertMemberAliasesTx(tx *sql.Tx, ctx context.Context, member_id int, alias string) error {
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

// scanRowsToMembers takes a pointer to *sql.Rows and scans those values into a slice of Members
func (m *postgresDBRepo) scanRowsToMembers(rows *sql.Rows) ([]models.Member, error) {
	var members []models.Member

	for rows.Next() {
		newMember := models.Member{}
		err := rows.Scan(&newMember.ID, &newMember.Name, &newMember.Email, &newMember.Active,
			&newMember.CreatedAt, &newMember.UpdatedAt)
		if err != nil {
			return members, err
		}

		// get aliases
		newMember.Aliases, err = m.getAliasesByMemberID(newMember.ID)
		if err != nil {
			return members, err
		}

		members = append(members, newMember)
	}
	err := rows.Err()
	if err != nil {
		return members, err
	}

	return members, nil
}

// scanRowsToMembersAliases takes a pointer to *sql.Rows and scans those values into a slice of MemberAliases
func scanRowsToMemberAliases(rows *sql.Rows) ([]models.MemberAlias, error) {
	var aliases []models.MemberAlias
	var memberId int

	for rows.Next() {
		m := models.MemberAlias{}
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

// AllMembers returns a slice of all members in database. Does not populate member aliases
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
	return m.scanRowsToMembers(rows)
}

// GetMemberByActive returns a slice of all members that have status = active. Does not populate member aliases
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
	return m.scanRowsToMembers(rows)
}

// GetMemberByID returns one member from a given id, populates aliases
func (m *postgresDBRepo) GetMemberByID(id int) (models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var v models.Member

	q := fmt.Sprintf(`SELECT id, %s FROM members
		WHERE id = $1`, memberCols)

	// execute our DB query
	row := m.DB.QueryRowContext(ctx, q, id)

	// scan single db row into member model
	err := row.Scan(&v.ID, &v.Name, &v.Email, &v.Active,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return v, err
	}

	// get member aliases
	v.Aliases, err = m.getAliasesByMemberID(id)
	if err != nil {
		return v, err
	}
	// q = fmt.Sprintf(`SELECT id, %s FROM member_aliases WHERE member_id = $1`, aliasCols)
	// rows, err := m.DB.QueryContext(ctx, q, id)
	// if err != nil {
	// 	return v, err
	// }
	// defer rows.Close()

	// v.Aliases, err = scanRowsToMemberAliases(rows)
	// if err != nil {
	// 	return v, err
	// }

	return v, nil
}

func (m *postgresDBRepo) getAliasesByMemberID(id int) ([]models.MemberAlias, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, %s FROM member_aliases WHERE member_id = $1`, aliasCols)
	rows, err := m.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMemberAliases(rows)
}

// UpdateMember updates a member in the database
func (m *postgresDBRepo) UpdateMember(v models.Member) error {
	return runInTx(m.DB, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		q := `UPDATE members SET
			name = $1,
			email = $2,
			is_active = $3,
			updated_at = $4
		WHERE id =  $5 `

		_, err := tx.ExecContext(ctx, q,
			v.Name,
			v.Email,
			v.Active,
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
		for _, a := range v.Aliases {
			err := insertMemberAliasesTx(tx, ctx, v.ID, a.Name)

			if err != nil {
				return err
			}

		}

		return nil
	})
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

// DeleteMember deletes one member by id. Also deletes all member_aliases with that member id
// We should almost never actually delete a member
// because it will be referenced by a lot of mileage logs.
func (m *postgresDBRepo) DeleteMember(id int) error {
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

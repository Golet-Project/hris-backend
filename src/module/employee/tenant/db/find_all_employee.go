package db

import (
	"context"
	"hris/module/shared/postgres"
	"hris/module/shared/primitive"
	"time"

	"github.com/jackc/pgx/v5"
)

type FindAllEmployeeOut struct {
	UID            string
	Email          string
	FullName       string
	BirthDate      time.Time
	Gender         primitive.Gender
	EmployeeStatus primitive.EmployeeStatus
	JoinDate       time.Time
	EndDate        primitive.Date
}

func (d *Db) FindAllEmployee(ctx context.Context, domain string) (out []FindAllEmployeeOut, err error) {
	if domain == "" {
		return nil, pgx.ErrNoRows
	}

	sql := `
	SELECT
		e.uid, e.email, e.full_name, e.birth_date, e.gender, e.employee_status,
		e.join_date, e.end_date
	FROM
		employee AS e`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row FindAllEmployeeOut

		err = rows.Scan(&row.UID, &row.Email, &row.FullName, &row.BirthDate, &row.Gender, &row.EmployeeStatus,
			&row.JoinDate, &row.EndDate,
		)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

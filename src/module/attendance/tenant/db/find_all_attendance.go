package db

import (
	"context"
	"hroost/module/shared/postgres"
	"hroost/module/shared/primitive"

	"github.com/jackc/pgx/v5"
)

type FindAllAttendanceOut struct {
	UID          string
	FullName     string
	CheckinTime  primitive.Time
	CheckoutTime primitive.Time
	ApprovedAt   primitive.Time
	ApprovedBy   primitive.String
}

func (d *Db) FindAllAttendance(ctx context.Context, domain string) (out []FindAllAttendanceOut, err error) {
	if domain == "" {
		return nil, pgx.ErrNoRows
	}

	var sql = `
	SELECT
		a.uid, e.full_name, a.checkin_time, a.checkout_time, a.approved_at,
		a.approved_by
	FROM
		attendance AS a
		INNER JOIN employee AS e ON e.uid = a.employee_uid`

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
		var row FindAllAttendanceOut

		err = rows.Scan(
			&row.UID, &row.FullName, &row.CheckinTime, &row.CheckoutTime, &row.ApprovedAt,
			&row.ApprovedBy,
		)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

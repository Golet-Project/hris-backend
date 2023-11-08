package db

import (
	"context"
	"hris/module/shared/postgres"
)

func (d *Db) CheckTodayAttendanceById(ctx context.Context, domain string, uid string) (exist bool, err error) {
	var sql = `
	SELECT
		COUNT (id)
	FROM
		attendance
	WHERE
		employee_uid = $1
		AND
		created_at::date = now()::date`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return false, err
	}

	var count int64
	err = conn.QueryRow(ctx, sql, uid).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

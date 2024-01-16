package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/module/shared/primitive"
)

func (d *Db) CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exist bool, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)

	var companyTimezoneSql = `
	SELECT
		timezone
	FROM
		tenant
	WHERE
		domain = $1`
	var companyTz int64
	err = masterConn.QueryRow(ctx, companyTimezoneSql, postgres.Domain(domain)).Scan(&companyTz)
	if err != nil {
		return
	}

	companyTimeNow, err := primitive.Timezone(companyTz).Now()
	if err != nil {
		return
	}
	startDate := companyTimeNow.Format("2006-01-02T00:00:00Z07:00")
	endDate := companyTimeNow.AddDate(0, 0, 1).Format("2006-01-02T00:00:00Z07:00")

	var sql = `
	SELECT
		COUNT(id)
	FROM
		attendance
	WHERE
		employee_uid = $1
		AND
		checkin_time >= $2
		AND
		checkin_time < $3`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return false, err
	}

	var count int64
	err = conn.QueryRow(ctx, sql, uid, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

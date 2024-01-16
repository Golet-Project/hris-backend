package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/module/shared/entities"
	"hroost/module/shared/primitive"

	"github.com/jackc/pgx/v5"
)

type GetTodayAttendanceIn struct {
	EmployeeUID string
	Timezone    primitive.Timezone
}

type GetTodayAttendanceOut struct {
	Timezone         primitive.Timezone
	AttendanceRadius primitive.Int64
	CheckinTime      primitive.Time
	CheckoutTime     primitive.Time
	ApprovedAt       primitive.Time

	StartWorkingHour primitive.Time
	EndWorkingHour   primitive.Time

	Company entities.Company
}

func (d *Db) GetTodayAttendance(ctx context.Context, domain string, param GetTodayAttendanceIn) (out GetTodayAttendanceOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	if domain == "" || !param.Timezone.Valid() || param.EmployeeUID == "" {
		return GetTodayAttendanceOut{}, pgx.ErrNoRows
	}

	var companySql = `
	SELECT
		latitude, longitude, attendance_radius
	FROM
		tenant
	WHERE
		domain = $1`

	err = masterConn.QueryRow(ctx, companySql, domain).Scan(&out.Company.Coordinate.Latitude, &out.Company.Coordinate.Longitude, &out.AttendanceRadius)
	if err != nil {
		return
	}

	var sql = `
	SELECT
		timezone, checkin_time, checkout_time, approved_at
	FROM
		attendance
	WHERE
		employee_uid = $1
		AND
		checkin_time::date = $2`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	now, err := param.Timezone.Now()
	if err != nil {
		return
	}
	todayDate := now.Format("2006-01-02")

	var tz int64
	err = conn.QueryRow(ctx, sql, param.EmployeeUID, todayDate).Scan(&tz, &out.CheckinTime, &out.CheckoutTime, &out.ApprovedAt)
	if err != nil {
		return
	}

	out.Timezone = primitive.Timezone(tz)

	return
}

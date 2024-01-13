package db

import (
	"context"
	"hroost/module/shared/postgres"
	"hroost/module/shared/primitive"
	"time"

	"github.com/jackc/pgx/v5"
)

type AddAttendanceIn struct {
	EmployeeUID string
	Timezone    primitive.Timezone
	Coordinate  primitive.Coordinate
}

func (d *Db) AddAttendance(ctx context.Context, domain string, in AddAttendanceIn) (err error) {
	if domain == "" || in.EmployeeUID == "" {
		return pgx.ErrNoRows
	}

	now, err := in.Timezone.Now()
	if err != nil {
		return
	}

	var sql = `
	INSERT INTO attendance
	(
		employee_uid, timezone, checkin_time, latitude, longitude
	) 
		VALUES 
	(
			$1, $2, $3, $4, $5
	)`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	_, err = conn.Exec(ctx, sql,
		in.EmployeeUID,
		in.Timezone.Value(),
		now.Format(time.RFC3339),
		in.Coordinate.Latitude,
		in.Coordinate.Longitude,
	)

	return
}

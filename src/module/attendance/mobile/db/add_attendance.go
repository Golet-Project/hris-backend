package db

import (
	"context"
	"hris/module/shared/postgres"
	"hris/module/shared/primitive"

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

	var sql = `
	INSERT INTO attendance
	(
		employee_id, timezone, latitude, longitude
	) 
		VALUES 
	(
			$1, $2, $3, $4
	)`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	_, err = conn.Exec(ctx, sql, in.EmployeeUID, in.Timezone.Value(), in.Coordinate.Latitude, in.Coordinate.Longitude)

	return
}

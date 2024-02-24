package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/attendance/model"
	"hroost/shared/primitive"
	"time"
)

func (d *Db) AddAttendance(ctx context.Context, domain string, in model.AddAttendanceIn) (repoError *primitive.RepoError) {
	if domain == "" || in.EmployeeUID == "" {
		return &primitive.RepoError{Issue: primitive.RepoErrorCodeDataNotFound}
	}

	now, err := in.Timezone.Now()
	if err != nil {
		return &primitive.RepoError{Issue: primitive.RepoErrorCodeServerError}
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

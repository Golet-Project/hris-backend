package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/attendance/model"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetTodayAttendance(ctx context.Context, domain string, param model.GetTodayAttendanceIn) (out model.GetTodayAttendanceOut, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	if domain == "" || !param.Timezone.Valid() || param.EmployeeUID == "" {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   err,
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	now, err := param.Timezone.Now()
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}
	todayDate := now.Format("2006-01-02")

	var tz int64
	err = conn.QueryRow(ctx, sql, param.EmployeeUID, todayDate).Scan(&tz, &out.CheckinTime, &out.CheckoutTime, &out.ApprovedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	out.Timezone = primitive.Timezone(tz)

	return
}

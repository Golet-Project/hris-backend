package db

import (
	"context"
	"errors"
	"fmt"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/homepage/model"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (d *Db) FindHomePage(ctx context.Context, domain string, param model.FindHomePageIn) (out model.FindHomePageOut, repoError *primitive.RepoError) {
	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	// get today attendance
	todayAttendance, err := getTodayAttendance(ctx, conn, param.UID, param.Timezone)
	if err != nil {
		return
	}

	out.TodayAttendance = todayAttendance

	return
}

func getTodayAttendance(ctx context.Context, conn *pgxpool.Pool, uid string, timezone primitive.Timezone) (out model.TodayAttendance, err error) {
	now, err := timezone.Now()
	if err != nil {
		return model.TodayAttendance{}, fmt.Errorf("error getting now time: %v", err)
	}

	todayDate := now.Format("2006-01-02")

	var sql = `
	SELECT
		timezone, checkin_time, checkout_time, approved_at
	FROM
		attendance
	WHERE
		employee_uid = $1
		AND
		checkin_time::date = $2`

	var tz int64
	err = conn.QueryRow(ctx, sql, uid, todayDate).Scan(&tz, &out.CheckinTime, &out.CheckoutTime, &out.ApprovedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.TodayAttendance{}, nil
		} else {
			return model.TodayAttendance{}, err
		}
	}

	out.Timezone = primitive.Timezone(tz)

	return
}

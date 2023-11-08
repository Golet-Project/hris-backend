package db

import (
	"context"
	"errors"
	"fmt"
	"hris/module/shared/postgres"
	"hris/module/shared/primitive"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FindHomePageIn struct {
	UID      string
	Timezone primitive.Timezone
}

type FindHomePageOut struct {
	TodayAttendance
}

func (d *Db) FindHomePage(ctx context.Context, domain string, param FindHomePageIn) (out FindHomePageOut, err error) {
	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	// get today attendance
	todayAttendance, err := getTodayAttendance(ctx, conn, param.UID, param.Timezone)
	if err != nil {
		return
	}

	out.TodayAttendance = todayAttendance

	return
}

type TodayAttendance struct {
	Timezone     primitive.Timezone
	CheckinTime  primitive.Time
	CheckoutTime primitive.Time
	ApprovedAt   primitive.Time
}

func getTodayAttendance(ctx context.Context, conn *pgxpool.Pool, uid string, timezone primitive.Timezone) (out TodayAttendance, err error) {
	now, err := timezone.Now()
	if err != nil {
		return TodayAttendance{}, fmt.Errorf("error getting now time: %v", err)
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
			return TodayAttendance{}, nil
		} else {
			return TodayAttendance{}, err
		}
	}

	out.Timezone = primitive.Timezone(tz)

	return
}

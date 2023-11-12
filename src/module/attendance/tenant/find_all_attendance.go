package tenant

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type FindAllAttendanceIn struct {
	Domain string
}

type Attendance struct {
	UID          string `json:"uid"`
	FullName     string `json:"full_name"`
	Date         string `json:"date"`
	CheckinTime  string `json:"checkin_time"`
	CheckoutTime string `json:"checkout_time"`
	ApprovedAt   string `json:"approved_at"`
	ApprovedBy   string `json:"approved_by"`
}

type FindAllAttendanceOut struct {
	primitive.CommonResult

	Attendance []Attendance
}

func (t *Tenant) FindAllAttendance(ctx context.Context, req FindAllAttendanceIn) (out FindAllAttendanceOut) {
	attendances, err := t.db.FindAllAttendance(ctx, req.Domain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	for _, attendance := range attendances {
		var a Attendance

		a.UID = attendance.UID
		a.FullName = attendance.FullName
		if attendance.CheckinTime.Valid {
			a.Date = attendance.CheckinTime.Time.Format("2006-01-02")
			a.CheckinTime = attendance.CheckinTime.Time.Format("15:04:05")
		}
		if attendance.CheckoutTime.Valid {
			a.CheckoutTime = attendance.CheckoutTime.Time.Format("15:04:05")
		}
		if attendance.ApprovedAt.Valid {
			a.ApprovedAt = attendance.ApprovedAt.Time.Format(time.RFC3339)
		}
		a.ApprovedBy = attendance.ApprovedBy.String

		out.Attendance = append(out.Attendance, a)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}

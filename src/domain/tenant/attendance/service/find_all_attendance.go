package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"hroost/module/shared/primitive"

	"github.com/jackc/pgx/v5"
)

type FindAllAttendanceIn struct {
	Domain string
}

type Attendance struct {
	UID          string `json:"uid"`
	FullName     string `json:"full_name"`
	CheckinTime  string `json:"checkin_time"`
	CheckoutTime string `json:"checkout_time"`
	ApprovedAt   string `json:"approved_at"`
	ApprovedBy   string `json:"approved_by"`
}

type FindAllAttendanceOut struct {
	primitive.CommonResult

	Attendance []Attendance
}

func (s *Service) FindAllAttendance(ctx context.Context, req FindAllAttendanceIn) (out FindAllAttendanceOut) {
	attendances, err := s.db.FindAllAttendance(ctx, req.Domain)
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
			a.CheckinTime = attendance.CheckinTime.Time.UTC().Format(time.RFC3339)
		}
		if attendance.CheckoutTime.Valid {
			a.CheckoutTime = attendance.CheckoutTime.Time.UTC().Format(time.RFC3339)
		}
		if attendance.ApprovedAt.Valid {
			a.ApprovedAt = attendance.ApprovedAt.Time.UTC().Format(time.RFC3339)
		}
		a.ApprovedBy = attendance.ApprovedBy.String

		out.Attendance = append(out.Attendance, a)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}

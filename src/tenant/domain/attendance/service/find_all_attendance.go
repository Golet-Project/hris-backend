package service

import (
	"context"
	"net/http"
	"time"

	"hroost/shared/primitive"
	"hroost/tenant/domain/attendance/model"
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

type FindAllAttendanceDb interface {
	FindAllAttendance(ctx context.Context, domain string) (out []model.FindAllAttendanceOut, err *primitive.RepoError)
}

type FindAllAttendance struct {
	Db FindAllAttendanceDb
}

func (s *FindAllAttendance) Exec(ctx context.Context, req FindAllAttendanceIn) (out FindAllAttendanceOut) {
	attendances, repoError := s.Db.FindAllAttendance(ctx, req.Domain)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
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

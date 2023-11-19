package mobile

import (
	"context"
	"errors"
	"hris/module/attendance/mobile/db"
	"hris/module/shared/entities"
	"hris/module/shared/primitive"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type GetTodayAttendanceIn struct {
	EmployeeUID string
	Timezone    primitive.Timezone
}

type GetTodayAttendanceOut struct {
	primitive.CommonResult

	Timezone         primitive.Timezone `json:"timezone"`
	CurrentTime      string             `json:"current_time"`
	CheckinTime      string             `json:"checkin_time"`
	CheckoutTime     string             `json:"checkout_time"`
	ApprovedAt       string             `json:"approved_at"`
	StartWorkingHour string             `json:"start_working_hour"`
	EndWorkingHour   string             `json:"end_working_hour"`
	AttendanceRadius int64              `json:"attendance_radius"` // radius in meters

	Company entities.Company `json:"company"`
}

func ValidateGetTodayAttendanceIn(req GetTodayAttendanceIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate timezone
	if !req.Timezone.Valid() {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeInvalidValue,
			Field:   "timezone",
			Message: "timezone header invalid",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

func (m *Mobile) GetTodayAttendance(ctx context.Context, req GetTodayAttendanceIn) (out GetTodayAttendanceOut) {
	// get the domain
	domain, err := m.userService.GetDomainByUid(ctx, req.EmployeeUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found", err)
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// validate request
	if err := ValidateGetTodayAttendanceIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "invalid request", err)
		return
	}

	todayAttendance, err := m.db.GetTodayAttendance(ctx, domain, db.GetTodayAttendanceIn{
		EmployeeUID: req.EmployeeUID,
		Timezone:    req.Timezone,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusInternalServerError, "error getting today attendance", err)
			return
		}
	}

	// build response
	out.Timezone = todayAttendance.Timezone
	out.CurrentTime = time.Now().UTC().Format(primitive.UtcRFC3339)
	if todayAttendance.CheckinTime.Valid {
		out.CheckinTime = todayAttendance.CheckinTime.Time.UTC().Format(primitive.UtcRFC3339)
	}
	if todayAttendance.CheckoutTime.Valid {
		out.CheckoutTime = todayAttendance.CheckoutTime.Time.UTC().Format(primitive.UtcRFC3339)
	}
	if todayAttendance.ApprovedAt.Valid {
		out.ApprovedAt = todayAttendance.ApprovedAt.Time.UTC().Format(primitive.UtcRFC3339)
	}
	st, err := time.Parse("15:04", "09:00")
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "error getting today attendance", err)
		return
	}
	out.StartWorkingHour = st.Format("15:04")

	et, err := time.Parse("15:04", "17:00")
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "error getting today attendance", err)
		return
	}
	out.EndWorkingHour = et.Format("15:04")

	out.AttendanceRadius = todayAttendance.AttendanceRadius.Int64
	out.Company = todayAttendance.Company
	if todayAttendance.Company.Address.Valid {
		out.Company.Address = todayAttendance.Company.Address
	} else {
		out.Company.Address = primitive.String{String: "", Valid: true}
	}

	out.SetResponse(http.StatusOK, "success")

	return
}

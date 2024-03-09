package service

import (
	"context"
	"hroost/mobile/domain/attendance/model"
	"hroost/shared/entities"
	"hroost/shared/primitive"
	"net/http"
	"time"

	"github.com/google/uuid"
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

type GetTodayAttendanceDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	GetTodayAttendance(ctx context.Context, domain string, in model.GetTodayAttendanceIn) (out model.GetTodayAttendanceOut, err *primitive.RepoError)
}

type GetTodayAttendance struct {
	Db     GetTodayAttendanceDb
	GetNow func() time.Time
}

func (s *GetTodayAttendance) Exec(ctx context.Context, req GetTodayAttendanceIn) (out GetTodayAttendanceOut) {
	// validate request
	if err := s.ValidateGetTodayAttendanceIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the domain
	domain, repoError := s.Db.GetDomainByUid(ctx, req.EmployeeUID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found", repoError)
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	todayAttendance, repoError := s.Db.GetTodayAttendance(ctx, domain, model.GetTodayAttendanceIn{
		EmployeeUID: req.EmployeeUID,
		Timezone:    req.Timezone,
	})
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeServerError:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	// build response
	out.Timezone = todayAttendance.Timezone
	out.CurrentTime = s.GetNow().UTC().Format(primitive.UtcRFC3339)
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

func (s *GetTodayAttendance) ValidateGetTodayAttendanceIn(req GetTodayAttendanceIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	if req.EmployeeUID == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "employee_uid",
			Message: "must not be empty",
		})
	} else {
		parsed, err := uuid.Parse(req.EmployeeUID)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_uid",
				Message: "must be a valid UUID",
			})
		}

		if parsed.Version() != 4 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_uid",
				Message: "must be a valid UUIDV4",
			})
		}
	}

	// validate timezone
	if req.Timezone == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "timezone",
			Message: "must not be empty",
		})
	} else {
		if !req.Timezone.Valid() {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "must be a valid timezone header",
			})
		}
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

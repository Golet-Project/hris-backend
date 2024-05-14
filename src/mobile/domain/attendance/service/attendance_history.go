package service

import (
	"context"
	"hroost/mobile/domain/attendance/model"
	"hroost/shared/primitive"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AttendanceHistoryIn struct {
	EmployeeId string

	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

type Attendance struct {
	ID               string `json:"id"`
	Date             string `json:"date"`
	CheckinTime      string `json:"checkin_time"`
	CheckoutTime     string `json:"checkout_time"`
	ApprovedAt       string `json:"approved_at"`
	AttendanceRadius int64  `json:"attendance_radius"` // radius in meters
}

type AttendanceHistoryOut struct {
	primitive.CommonResult

	Length int64 `json:"-"`

	Attendances []Attendance `json:"attendances"`
}

type AttendanceHistoryDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	FindAttendanceHistory(ctx context.Context, domain string, in model.FindAttendanceHistoryIn) (length int64, attendances []model.FindAttendanceHistoryOut, err *primitive.RepoError)
}

type AttendanceHistory struct {
	Db AttendanceHistoryDb

	GetNow func() time.Time
}

func (s *AttendanceHistory) Exec(ctx context.Context, req AttendanceHistoryIn) (out AttendanceHistoryOut) {
	if err := s.ValidateRequest(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed")
		return
	}

	domain, repoError := s.Db.GetDomainByUid(ctx, req.EmployeeId)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	if req.StartDate == "" && req.EndDate == "" {
		now := s.GetNow()
		req.StartDate = now.AddDate(0, 0, -now.Day()+1).Format("2006-01-02")
		req.EndDate = now.AddDate(0, 1, -now.Day()).Format("2006-01-02")
	}

	attendanceHistoryCount, attendanceHistories, repoError := s.Db.FindAttendanceHistory(ctx, domain, model.FindAttendanceHistoryIn{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusOK, "success")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	out.Length = attendanceHistoryCount
	out.Attendances = make([]Attendance, len(attendanceHistories))

	for i, attendanceHistory := range attendanceHistories {
		a := Attendance{
			ID:               attendanceHistory.ID,
			AttendanceRadius: attendanceHistory.Radius,
		}
		if attendanceHistory.Date.Valid {
			a.Date = attendanceHistory.Date.String
		}
		if attendanceHistory.CheckinTime.Valid {
			a.CheckinTime = attendanceHistory.CheckinTime.Time.Format(primitive.UtcRFC3339)
		}
		if attendanceHistory.CheckoutTime.Valid {
			a.CheckoutTime = attendanceHistory.CheckoutTime.Time.Format(primitive.UtcRFC3339)
		}
		if attendanceHistory.ApprovedAt.Valid {
			a.ApprovedAt = attendanceHistory.ApprovedAt.Time.Format(primitive.UtcRFC3339)
		}

		out.Attendances[i] = a
	}

	out.SetResponse(http.StatusOK, "success")
	return
}

func (s *AttendanceHistory) ValidateRequest(req AttendanceHistoryIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// employee_id
	if req.EmployeeId == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "employee_id",
			Message: "must not be empty",
		})
	} else {
		parsed, err := uuid.Parse(req.EmployeeId)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_id",
				Message: "must be a valid UUID",
			})
		}
		if parsed.Version() != 4 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_id",
				Message: "must be a valid UUIDV4",
			})
		}
	}

	// start_date
	if req.StartDate == "" {
		if req.EndDate != "" {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "start_date",
				Message: "must not be empty if end_date is not empty",
			})
		}
	} else {
		if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "start_date",
				Message: "must be in the YYYY-MM-DD format",
			})
		}
	}

	// end_date
	if req.EndDate == "" {
		if req.StartDate != "" {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "end_date",
				Message: "must not be empty if start_date is not empty",
			})
		}
	} else {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "end_date",
				Message: "must be in the YYYY-MM-DD format",
			})
		}
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{Issues: allIssues}
	}

	return nil
}

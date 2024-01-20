package service

import (
	"context"
	"errors"
	"hroost/mobile/domain/attendance/db"
	"hroost/shared/primitive"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AddAttendanceIn struct {
	UID        string
	Timezone   primitive.Timezone
	Coordinate primitive.Coordinate `json:"coordinate"`
}

type AddAttendanceOut struct {
	primitive.CommonResult
}

func ValidateAddAttendanceRequest(req AddAttendanceIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate uid
	if len(req.UID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "uid",
			Message: "uid is required",
		})
	} else {
		_, err := uuid.Parse(req.UID)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid uuid",
			})
		}
	}

	// validate timezone
	if int(req.Timezone) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "timezone",
			Message: "timezone is required",
		})
	} else {
		if !req.Timezone.Valid() {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "timezone is not valid",
			})
		}
	}

	// validate coordinate
	if req.Coordinate == (primitive.Coordinate{}) {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "coordinate",
			Message: "coordinate is required",
		})
	} else {
		// validate coordinate.latitude
		if req.Coordinate.Latitude < -90 || req.Coordinate.Latitude > 90 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "recipient.latitude",
				Message: "must be between -90 and 90",
			})
		}

		// validate coordinate.longitude
		if req.Coordinate.Longitude < -180 || req.Coordinate.Longitude > 180 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "recipient.longitude",
				Message: "must be between -180 and 180",
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

func (s *Service) AddAttendance(ctx context.Context, req AddAttendanceIn) (out AddAttendanceOut) {
	// validate request
	if err := ValidateAddAttendanceRequest(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the domain
	domain, err := s.userService.GetDomainByUid(ctx, req.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// validate user
	exist, err := s.db.CheckEmployeeById(ctx, domain, req.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found", err)
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "interal server error", err)
			return
		}
	}
	if !exist {
		out.SetResponse(http.StatusNotFound, "employee not found")
		return
	}

	// check if attendance already exist
	exists, err := s.db.CheckTodayAttendanceById(ctx, domain, req.UID, req.Timezone)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "interal server error", err)
		return
	}
	if exists {
		out.SetResponse(http.StatusConflict, "attendance already exist")
		return
	}

	// input attendance
	err = s.db.AddAttendance(ctx, domain, db.AddAttendanceIn{
		EmployeeUID: req.UID,
		Timezone:    req.Timezone,
		Coordinate: primitive.Coordinate{
			Latitude:  req.Coordinate.Latitude,
			Longitude: req.Coordinate.Longitude,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	out.SetResponse(http.StatusCreated, "success")

	return
}

package service

import (
	"context"
	"errors"
	"hroost/mobile/domain/attendance/model"
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

type AddAttendanceDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	// TODO: move into separate struct
	EmployeeExistsById(ctx context.Context, domain string, uid string) (exists bool, err *primitive.RepoError)
	CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exists bool, err *primitive.RepoError)

	AddAttendance(ctx context.Context, domain string, data model.AddAttendanceIn) (err *primitive.RepoError)
}

type AddAttendance struct {
	Db AddAttendanceDb
}

func (s *AddAttendance) Exec(ctx context.Context, req AddAttendanceIn) (out AddAttendanceOut) {
	// validate request
	if err := s.ValidateAddAttendanceRequest(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the domain
	domain, repoError := s.Db.GetDomainByUid(ctx, req.UID)
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

	// validate user
	exist, repoError := s.Db.EmployeeExistsById(ctx, domain, req.UID)
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

	if !exist {
		out.SetResponse(http.StatusNotFound, "employee not found")
		return
	}

	// check if attendance already exist
	exists, err := s.Db.CheckTodayAttendanceById(ctx, domain, req.UID, req.Timezone)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if exists {
		out.SetResponse(http.StatusConflict, "attendance already exist")
		return
	}

	// input attendance
	err = s.Db.AddAttendance(ctx, domain, model.AddAttendanceIn{
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

func (s *AddAttendance) ValidateAddAttendanceRequest(req AddAttendanceIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate uid
	if len(req.UID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "uid",
			Message: "uid is required",
		})
	} else {
		parsed, err := uuid.Parse(req.UID)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid UUID",
			})
		}

		if parsed.Version() != 4 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid UUIDV4",
			})
		}
	}

	// validate timezone
	if int(req.Timezone) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "timezone",
			Message: "timezone header is required",
		})
	} else {
		if !req.Timezone.Valid() {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "timezone header has an invalid value",
			})
		}
	}

	// validate coordinate
	if req.Coordinate == (primitive.Coordinate{}) {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "coordinate",
			Message: "must not empty",
		})
	} else {
		// validate coordinate.latitude
		if req.Coordinate.Latitude < -90 || req.Coordinate.Latitude > 90 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.latitude",
				Message: "must be between -90 and 90",
			})
		}

		// validate coordinate.longitude
		if req.Coordinate.Longitude < -180 || req.Coordinate.Longitude > 180 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.longitude",
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

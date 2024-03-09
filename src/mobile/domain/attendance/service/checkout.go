package service

import (
	"context"
	"hroost/shared/primitive"
	"net/http"

	"github.com/google/uuid"
)

type CheckoutIn struct {
	UID      string
	Timezone primitive.Timezone
}

type CheckoutOut struct {
	primitive.CommonResult
}

type CheckoutDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exists bool, err *primitive.RepoError)
	Checkout(ctx context.Context, domaing string, uid string) (rowsAffected int64, err *primitive.RepoError)
}

type Checkout struct {
	Db CheckoutDb
}

func (s *Checkout) Exec(ctx context.Context, req CheckoutIn) (out CheckoutOut) {
	// validate request payload
	err := s.ValidateRequestPayload(req)
	if err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the domain
	domain, repoError := s.Db.GetDomainByUid(ctx, req.UID)
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

	// validate checkin
	attendanceInExists, repoError := s.Db.CheckTodayAttendanceById(ctx, domain, req.UID, req.Timezone)
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

	if !attendanceInExists {
		out.SetResponse(http.StatusBadRequest, "you haven't checkin yet")
		return
	}

	// insert checkout time
	rowsAffected, repoError := s.Db.Checkout(ctx, domain, req.UID)
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

	if rowsAffected == 0 {
		out.SetResponse(http.StatusNotFound, "you already checkout")
		return
	}

	out.SetResponse(http.StatusCreated, "success")

	return
}

func (s *Checkout) ValidateRequestPayload(in CheckoutIn) (err *primitive.RequestValidationError) {
	var issues []primitive.RequestValidationIssue
	// uid
	if in.UID == "" {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "uid",
			Message: "must not be empty",
		})
	} else {
		parsed, err := uuid.Parse(in.UID)
		if err != nil {
			issues = append(issues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "must be a valid UUID",
			})
		}

		if parsed.Version() != 4 {
			issues = append(issues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "must be a valid UUIDV4",
			})
		}
	}

	// timezone
	if in.Timezone == 0 {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "timezone",
			Message: "must not be empty",
		})
	} else {
		if !in.Timezone.Valid() {
			issues = append(issues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "must be a valid timezone",
			})
		}
	}

	if len(issues) > 0 {
		err = &primitive.RequestValidationError{
			Issues: issues,
		}
	}

	return err
}

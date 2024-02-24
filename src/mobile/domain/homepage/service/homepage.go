package service

import (
	"context"
	"hroost/mobile/domain/homepage/model"
	"hroost/shared/primitive"
	"net/http"
	"time"
)

type HomePageIn struct {
	UID      string
	Timezone primitive.Timezone
}

type TodayAttendance struct {
	Timezone     primitive.Timezone `json:"timezone"`
	CheckinTime  string             `json:"checkin_time"`
	CheckoutTime string             `json:"checkout_time"`
	ApprovedAt   string             `json:"approved_at"`
}

type HomePageOut struct {
	primitive.CommonResult

	TodayAttendance `json:"today_attendance"`
}

type HomePageDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	FindHomePage(ctx context.Context, domain string, query model.FindHomePageIn) (out model.FindHomePageOut, err *primitive.RepoError)
}

type HomePage struct {
	Db HomePageDb
}

func (s *HomePage) Exec(ctx context.Context, req HomePageIn) (out HomePageOut) {
	// validate request
	if err := s.ValidateHomePageIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "invalid request", err)
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
			out.SetResponse(http.StatusInternalServerError, "internal server eror", repoError)
			return
		}
	}

	// get the homepage data
	homepageData, repoError := s.Db.FindHomePage(ctx, domain, model.FindHomePageIn{
		UID:      req.UID,
		Timezone: req.Timezone,
	})
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "homepage data not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server eror", repoError)
			return
		}
	}

	out.TodayAttendance = TodayAttendance{
		Timezone: homepageData.Timezone,
	}

	if homepageData.CheckinTime.Valid {
		out.TodayAttendance.CheckinTime = homepageData.CheckinTime.Time.Format(time.RFC3339)
	}
	if homepageData.CheckoutTime.Valid {
		out.TodayAttendance.CheckoutTime = homepageData.CheckoutTime.Time.Format(time.RFC3339)
	}
	if homepageData.ApprovedAt.Valid {
		out.TodayAttendance.ApprovedAt = homepageData.ApprovedAt.Time.Format(time.RFC3339)
	}

	out.SetResponse(http.StatusOK, "success")

	return
}

func (s *HomePage) ValidateHomePageIn(req HomePageIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

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

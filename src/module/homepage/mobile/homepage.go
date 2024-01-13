package mobile

import (
	"context"
	"errors"
	"hroost/module/homepage/mobile/db"
	"hroost/module/shared/primitive"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
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

func ValidateHomePageIn(req HomePageIn) *primitive.RequestValidationError {
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

func (m *Mobile) HomePage(ctx context.Context, req HomePageIn) (out HomePageOut) {
	// validate request
	if err := ValidateHomePageIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "invalid request", err)
		return
	}

	// get the domain
	domain, err := m.userService.GetDomainByUid(ctx, req.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server eror", err)
			return
		}
	}

	// get the homepage data
	homepageData, err := m.db.FindHomePage(ctx, domain, db.FindHomePageIn{
		UID:      req.UID,
		Timezone: req.Timezone,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "homepage data not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server eror", err)
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

package service

import (
	"context"
	"errors"
	"hroost/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type CheckoutIn struct {
	UID      string
	Timezone primitive.Timezone
}

type CheckoutOut struct {
	primitive.CommonResult
}

func (s *Service) Checkout(ctx context.Context, req CheckoutIn) (out CheckoutOut) {
	// get the domain
	domain, err := s.userService.GetDomainByUid(ctx, req.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// validate checkin
	attendanceInExists, err := s.db.CheckTodayAttendanceById(ctx, domain, req.UID, req.Timezone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}
	if !attendanceInExists {
		out.SetResponse(http.StatusBadRequest, "you haven't checkin yet")
		return
	}

	// insert checkout time
	rowsAffected, err := s.db.Checkout(ctx, domain, req.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	} else {
		if rowsAffected == 0 {
			out.SetResponse(http.StatusNotFound, "you already checkout")
			return
		}
	}

	out.SetResponse(http.StatusCreated, "success checkout")

	return
}

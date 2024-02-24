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

type CheckoutDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exists bool, err *primitive.RepoError)
	Checkout(ctx context.Context, domaing string, uid string) (rowsAffected int64, err *primitive.RepoError)
}

type Checkout struct {
	Db CheckoutDb
}

func (s *Checkout) Exec(ctx context.Context, req CheckoutIn) (out CheckoutOut) {
	// get the domain
	domain, err := s.Db.GetDomainByUid(ctx, req.UID)
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
	attendanceInExists, err := s.Db.CheckTodayAttendanceById(ctx, domain, req.UID, req.Timezone)
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
	rowsAffected, err := s.Db.Checkout(ctx, domain, req.UID)
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

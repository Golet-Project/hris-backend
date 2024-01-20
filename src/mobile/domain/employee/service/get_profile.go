package service

import (
	"context"
	"errors"
	"hroost/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Employee struct {
	UID            string           `json:"uid"`
	FullName       string           `json:"full_name"`
	Email          string           `json:"email"`
	Gender         primitive.Gender `json:"gender"`
	BirthDate      primitive.Date   `json:"birth_date"`
	JoinDate       primitive.Date   `json:"join_date"`
	ProfilePicture primitive.String `json:"profile_picture"`
	Address        primitive.String `json:"address"`
}

type GetProfileOut struct {
	primitive.CommonResult

	Employee Employee `json:"employee"`
}

type GetProfileIn struct {
	UID string
}

func (s *Service) GetProfile(ctx context.Context, in GetProfileIn) (out GetProfileOut) {
	// get user domain
	domain, err := s.userService.GetDomainByUid(ctx, in.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found", err)
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	employee, err := s.db.GetEmployeeDetail(ctx, domain, in.UID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	out.Employee = Employee{
		UID:            employee.UID,
		FullName:       employee.FullName,
		Email:          employee.Email,
		Gender:         employee.Gender,
		BirthDate:      employee.BirthDate,
		JoinDate:       employee.JoinDate,
		ProfilePicture: employee.ProfilePicture,
		Address:        employee.Address,
	}
	out.SetResponse(http.StatusOK, "success")
	return
}

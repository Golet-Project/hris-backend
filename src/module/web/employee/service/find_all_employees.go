package service

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type FindAllEmployee struct {
	UID          string    `json:"uid"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	ProvinceName string    `json:"province_name"`
	RegencyName  string    `json:"regency_name"`
	DistrictName string    `json:"district_name"`
	VillageName  string    `json:"village_name"`
	RegisteredAt time.Time `json:"registered_at"`
}

type FindAllEmployeeOut struct {
	primitive.CommonResult

	Data []FindAllEmployee
}

func (s EmployeeService) FindAllEmployee(ctx context.Context) (out FindAllEmployeeOut) {
	// get all employee
	employees, err := s.EmployeeRepo.FindAllEmployee(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "employees data")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	for _, employee := range employees {
		out.Data = append(out.Data, FindAllEmployee{
			UID:          employee.UID,
			Email:        employee.Email,
			FullName:     employee.FullName,
			ProvinceName: employee.ProvinceName,
			RegencyName:  employee.RegencyName,
			VillageName:  employee.VillageName,
			RegisteredAt: employee.RegisteredAt,
		})
	}

	out.SetResponse(http.StatusOK, "employess data")

	return out
}

package service

import (
	"context"
	"hroost/shared/primitive"
	"hroost/tenant/domain/employee/model"
	"net/http"
)

type CreateEmployeeIn struct {
	Domain         string
	Email          string                   `json:"email"`
	FirstName      string                   `json:"first_name"`
	LastName       string                   `json:"last_name"`
	Gender         primitive.Gender         `json:"gender"`
	BirthDate      string                   `json:"birth_date"`
	Address        string                   `json:"address"`
	ProvinceId     string                   `json:"province_id"`
	RegencyId      string                   `json:"regency_id"`
	DistrictId     string                   `json:"district_id"`
	VillageId      string                   `json:"village_id"`
	JoinDate       string                   `json:"join_date"`
	EmployeeStatus primitive.EmployeeStatus `json:"employee_status"`
}

type CreateEmployeeOut struct {
	primitive.CommonResult
}

type CreateEmployeeDb interface {
	CreateEmployee(ctx context.Context, data model.CreateEmployeeIn) (err *primitive.RepoError)
}

type CreateEmployee struct {
	Db CreateEmployeeDb
}

func (s *CreateEmployee) Exec(ctx context.Context, req CreateEmployeeIn) (out CreateEmployeeOut) {
	// TODO: generate initial password and send to user's email
	err := s.Db.CreateEmployee(ctx, model.CreateEmployeeIn{

		Domain:         req.Domain,
		Email:          req.Email,
		Password:       "todo",
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Gender:         req.Gender,
		BirthDate:      req.BirthDate,
		Address:        req.Address,
		ProvinceId:     req.ProvinceId,
		RegencyId:      req.RegencyId,
		DistrictId:     req.DistrictId,
		VillageId:      req.VillageId,
		JoinDate:       req.JoinDate,
		EmployeeStatus: req.EmployeeStatus,
	})
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusCreated, "success")
	return
}

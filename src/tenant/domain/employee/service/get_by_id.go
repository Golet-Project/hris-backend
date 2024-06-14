package service

import (
	"context"
	"hroost/shared/primitive"
	"hroost/tenant/domain/employee/model"
	"net/http"
)

type GetByIdIn struct {
	Id string `params:"id"`
}

type GetByIdOut struct {
	primitive.CommonResult
	Address struct {
		Detail       string `json:"detail"`
		ProvinceId   string `json:"province_id"`
		ProvinceName string `json:"province_name"`
		RegencyId    string `json:"regency_id"`
		RegencyName  string `json:"regency_name"`
		DistrictId   string `json:"district_id"`
		DistrictName string `json:"district_name"`
		VillageId    string `json:"village_id"`
		VillageName  string `json:"village_name"`
	} `json:"address"`
	Id             string                   `json:"id"`
	Email          string                   `json:"email"`
	FullName       string                   `json:"full_name"`
	Gender         string                   `json:"gender"`
	EmployeeStatus primitive.EmployeeStatus `json:"employee_status"`
	BirthDate      primitive.Date           `json:"birth_date"`
	JoinDate       primitive.Date           `json:"join_date"`
}

type GetByIdDb interface {
	GetDomainById(ctx context.Context, id string) (domain string, err *primitive.RepoError)
	GetById(ctx context.Context, domain, id string) (out model.GetByIdOut, err *primitive.RepoError)
}

type GetById struct {
	Db GetByIdDb
}

func (s *GetById) Exec(ctx context.Context, req GetByIdIn) (out GetByIdOut) {
	if req.Id == "" {
		out.SetResponse(http.StatusNotFound, "employee not found", nil)
		return
	}

	domain, repoError := s.Db.GetDomainById(ctx, req.Id)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found", nil)
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	employee, repoError := s.Db.GetById(ctx, domain, req.Id)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found", nil)
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	out.Id = employee.Id
	out.Email = employee.Email
	out.FullName = employee.FullName
	out.Gender = employee.Gender
	out.EmployeeStatus = employee.EmployeeStatus
	out.BirthDate = employee.BirthDate
	out.JoinDate = employee.JoinDate
	out.Address.Detail = employee.Address.Detail.String
	out.Address.ProvinceId = employee.Address.ProvinceId.String
	out.Address.ProvinceName = employee.Address.ProvinceName.String
	out.Address.RegencyId = employee.Address.RegencyId.String
	out.Address.RegencyName = employee.Address.RegencyName.String
	out.Address.DistrictId = employee.Address.DistrictId.String
	out.Address.DistrictName = employee.Address.DistrictName.String
	out.Address.VillageId = employee.Address.VillageId.String
	out.Address.VillageName = employee.Address.VillageName.String

	out.SetResponse(http.StatusOK, "success")

	return
}

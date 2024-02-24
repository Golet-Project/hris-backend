package service

import (
	"context"
	"hroost/mobile/domain/employee/model"
	"hroost/shared/primitive"
	"net/http"
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

type GetProfileDb interface {
	GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError)
	GetEmployeeDetail(ctx context.Context, domain string, uid string) (employee model.GetEmployeeDetailOut, err *primitive.RepoError)
}

type GetProfile struct {
	Db GetProfileDb
}

func (s *GetProfile) Exec(ctx context.Context, in GetProfileIn) (out GetProfileOut) {
	// get user domain
	domain, repoError := s.Db.GetDomainByUid(ctx, in.UID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "user not found", repoError)
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	employee, repoError := s.Db.GetEmployeeDetail(ctx, domain, in.UID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found", repoError)
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
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

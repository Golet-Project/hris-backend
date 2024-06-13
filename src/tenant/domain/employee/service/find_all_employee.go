package service

import (
	"context"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"hroost/tenant/domain/employee/model"
	"net/http"
)

type FindAllEmployeeIn struct {
	Domain string
}

type Employee struct {
	Id             string                   `json:"id"`
	FullName       string                   `json:"full_name"`
	BirthDate      string                   `json:"birth_date"`
	Gender         primitive.Gender         `json:"gender"`
	Email          string                   `json:"email"`
	PhoneNumber    string                   `json:"phone_number"`
	JoinDate       string                   `json:"join_date"`
	EmployeeStatus primitive.EmployeeStatus `json:"employee_status"`
	EndDate        primitive.Date           `json:"end_date"`
	Age            int                      `json:"age"`
}

type FindAllEmployeeOut struct {
	primitive.CommonResult

	Employee []Employee `json:"employee"`
}

type FindAllEmployeeDb interface {
	FindAllEmployee(ctx context.Context, domain string) (out []model.FindAllEmployeeOut, err *primitive.RepoError)
}

type FindAllEmployee struct {
	Db FindAllEmployeeDb
}

// FindAllEmployee find all employee
func (s *FindAllEmployee) Exec(ctx context.Context, req FindAllEmployeeIn) (out FindAllEmployeeOut) {
	// validate the request
	if err := s.ValidateFindAllEmployeeIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed")
		return
	}

	// find the employee
	employees, repoError := s.Db.FindAllEmployee(ctx, req.Domain)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusOK, "success")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	// map the responses
	s.mapFindAllEmployee(employees, &out)

	out.SetResponse(http.StatusOK, "success")
	return
}

// mapFindALlEmployee map the data returned from database into response
func (s *FindAllEmployee) mapFindAllEmployee(in []model.FindAllEmployeeOut, out *FindAllEmployeeOut) {
	for _, employee := range in {
		var o Employee
		o.Id = employee.Id
		o.FullName = employee.FullName
		o.BirthDate = employee.BirthDate.Format("2006-01-02")
		o.Gender = employee.Gender
		o.Age = utils.CalculateAge(employee.BirthDate)
		o.Email = employee.Email
		o.PhoneNumber = ""
		o.JoinDate = employee.JoinDate.Format("2006-01-02")
		o.EndDate = employee.EndDate
		o.EmployeeStatus = employee.EmployeeStatus

		out.Employee = append(out.Employee, o)
	}
}

func (s *FindAllEmployee) ValidateFindAllEmployeeIn(req FindAllEmployeeIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	if req.Domain == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "domain",
			Message: "domain is required",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

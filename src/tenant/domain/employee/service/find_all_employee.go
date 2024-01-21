package service

import (
	"context"
	"errors"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"hroost/tenant/domain/employee/db"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type FindAllEmployeeIn struct {
	Domain string
}

type FindAllEmployee struct {
	UID            string                   `json:"uid"`
	FullName       string                   `json:"full_name"`
	Gender         primitive.Gender         `json:"gender"`
	Age            int                      `json:"age"`
	Email          string                   `json:"email"`
	PhoneNumber    string                   `json:"phone_number"`
	JoinDate       string                   `json:"join_date"`
	EndDate        primitive.Date           `json:"end_date"`
	EmployeeStatus primitive.EmployeeStatus `json:"employee_status"`
}

type FindAllEmployeeOut struct {
	primitive.CommonResult

	Employee []FindAllEmployee `json:"employee"`
}

func ValidateFindAllEmployeeIn(req FindAllEmployeeIn) *primitive.RequestValidationError {
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

// FindAllEmployee find all employee
func (s Service) FindAllEmployee(ctx context.Context, req FindAllEmployeeIn) (out FindAllEmployeeOut) {
	// validate the request
	if err := ValidateFindAllEmployeeIn(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed")
		return
	}

	// find the employee
	employees, err := s.db.FindAllEmployee(ctx, req.Domain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusOK, "success")
			return
		} else {
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
func (s *Service) mapFindAllEmployee(in []db.FindAllEmployeeOut, out *FindAllEmployeeOut) {
	for _, employee := range in {
		var o FindAllEmployee
		o.UID = employee.UID
		o.FullName = employee.FullName
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

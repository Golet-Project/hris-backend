package service

import (
	"context"
	"errors"
	employeeRepo "hris/module/employee/repo/employee"
	"hris/module/shared/primitive"
	"hris/module/shared/utils"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type WebFindAllEmployeesIn struct {
}

type WebFindAllEmployee struct {
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

type WebFindAllEmployeesOut struct {
	primitive.CommonResult
	Employees []WebFindAllEmployee
}

func (s *WebEmployeeService) FindAllEmployees(ctx context.Context) (out WebFindAllEmployeesOut) {
	// find the employee
	employees, err := s.EmployeeRepo.WebFindAllEmployees(ctx)
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
	s.mapFindAllEmployees(employees, &out)

	out.SetResponse(http.StatusOK, "success")
	return
}

func (s *WebEmployeeService) mapFindAllEmployees(in []employeeRepo.WebFindAllEmployeesOut, out *WebFindAllEmployeesOut) {
	for _, employee := range in {
		var o WebFindAllEmployee
		o.UID = employee.UID
		o.FullName = employee.FullName
		o.Gender = employee.Gender
		o.Age = utils.CalculateAge(employee.BirthDate)
		o.Email = employee.Email
		o.PhoneNumber = ""
		o.JoinDate = employee.JoinDate.Format("2006-01-02")
		o.EndDate = employee.EndDate
		o.EmployeeStatus = employee.EmployeeStatus

		out.Employees = append(out.Employees, o)
	}
}

package tenant

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"hris/module/shared/utils"
	"net/http"

	"hris/module/employee/tenant/db"

	"github.com/jackc/pgx/v5"
)

type FindAllEmployeeIn struct {
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

// FindAllEmployee find all employee
func (t *Tenant) FindAllEmployee(ctx context.Context) (out FindAllEmployeeOut) {
	// find the employee
	employees, err := t.db.FindAllEmployee(ctx)
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
	t.mapFindAllEmployee(employees, &out)

	out.SetResponse(http.StatusOK, "success")
	return
}

// mapFindALlEmployee map the data returned from database into response
func (t *Tenant) mapFindAllEmployee(in []db.FindAllEmployeeOut, out *FindAllEmployeeOut) {
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
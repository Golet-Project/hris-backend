package employee

import (
	"hris/module/web/employee/presentation/rest"
	"hris/module/web/employee/repo/employee"
	"hris/module/web/employee/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB *pgxpool.Pool
}

type Employee struct {
	EmployeePresenter *rest.EmployeePresenter
}

func InitEmployee(d *Dependency) *Employee {
	employeeRepo := &employee.Repository{
		DB: d.DB,
	}

	employeeService := service.NewEmployeeService(employeeRepo)

	return &Employee{
		EmployeePresenter: &rest.EmployeePresenter{
			EmployeeService: employeeService,
		},
	}
}

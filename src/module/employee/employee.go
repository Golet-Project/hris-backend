package employee

import (
	"hris/module/employee/presenter/rest"
	"hris/module/employee/repo/employee"
	"hris/module/employee/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Employee struct {
	EmployeePresenter *rest.EmployeePresenter
}

type Dependency struct {
	DB *pgxpool.Pool
}

func InitEmployee(d *Dependency) *Employee {
	if d.DB == nil {
		log.Fatal("[x] Employee package require a database connection")
	}

	employeeRepo := employee.Repository{
		DB: d.DB,
	}

	webEmployeeService := service.NewWebEmployeeService(&employeeRepo)

	return &Employee{
		EmployeePresenter: &rest.EmployeePresenter{
			WebAuthService: webEmployeeService,
		},
	}
}

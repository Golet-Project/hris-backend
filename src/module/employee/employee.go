package employee

import (
	"hris/module/employee/presenter/rest"
	"hris/module/employee/tenant"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Employee struct {
	EmployeePresentation *rest.EmployeePresentation
}

type Dependency struct {
	DB *pgxpool.Pool
}

func InitEmployee(d *Dependency) *Employee {
	if d.DB == nil {
		log.Fatal("[x] Employee package require a database connection")
	}

	tenantEmployeeService := tenant.New(&tenant.Dependency{
		Pg: d.DB,
	})

	return &Employee{
		EmployeePresentation: &rest.EmployeePresentation{
			Tenant: tenantEmployeeService,
		},
	}
}

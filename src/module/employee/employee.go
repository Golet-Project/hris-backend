package employee

import (
	"hris/module/employee/presentation/rest"
	"hris/module/employee/tenant"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Employee struct {
	EmployeePresentation *rest.EmployeePresentation
}

type Dependency struct {
	MasterDB *pgxpool.Pool
}

func InitEmployee(d *Dependency) *Employee {
	if d.MasterDB == nil {
		log.Fatal("[x] Employee package require a database connection")
	}

	tenantEmployeeService := tenant.New(&tenant.Dependency{
		MasterConn: d.MasterDB,
	})

	employeePresentation := rest.New(tenantEmployeeService)

	return &Employee{
		EmployeePresentation: employeePresentation,
	}
}

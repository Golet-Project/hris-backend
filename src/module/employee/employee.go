package employee

import (
	"hroost/module/employee/mobile"
	"hroost/module/employee/presentation/rest"
	"hroost/module/employee/tenant"
	"hroost/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	userService "hroost/module/user/service"
)

type Employee struct {
	EmployeePresentation *rest.EmployeePresentation
}

type Dependency struct {
	MasterDB   *pgxpool.Pool
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func InitEmployee(d *Dependency) *Employee {
	if d.MasterDB == nil {
		log.Fatal("[x] Employee package require a database connection")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Employee package require a database resolver")
	}
	if d.UserService == nil {
		log.Fatal("[x] Employee package require a user service")
	}

	tenantEmployeeService := tenant.New(&tenant.Dependency{
		MasterConn: d.MasterDB,
		PgResolver: d.PgResolver,
	})
	mobileEmployeeService := mobile.New(mobile.Dependency{
		MasterConn: d.MasterDB,
		PgResolver: d.PgResolver,

		UserService: d.UserService,
	})

	employeePresentation := rest.New(tenantEmployeeService, mobileEmployeeService)

	return &Employee{
		EmployeePresentation: employeePresentation,
	}
}

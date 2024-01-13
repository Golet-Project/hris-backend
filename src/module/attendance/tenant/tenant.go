package tenant

import (
	"hroost/module/attendance/tenant/db"
	"hroost/module/shared/postgres"
	userService "hroost/module/user/service"
	"log"
)

type Tenant struct {
	db *db.Db

	userService *userService.Service
}

type Dependency struct {
	PgResolver *postgres.Resolver

	UserService *userService.Service
}

func New(d *Dependency) *Tenant {
	if d.PgResolver == nil {
		log.Fatal("[x] postgres resolver is required on tenant module")
	}
	if d.UserService == nil {
		log.Fatal("[x] user service is required on tenant module")
	}

	return &Tenant{
		db: db.New(&db.Dependency{
			PgResolver: d.PgResolver,
		}),

		userService: d.UserService,
	}
}

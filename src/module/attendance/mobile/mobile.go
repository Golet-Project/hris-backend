package mobile

import (
	"hris/module/attendance/mobile/db"
	"hris/module/shared/postgres"
	"log"

	userService "hris/module/user/service"
)

type Mobile struct {
	db *db.Db

	// other service
	userService *userService.Service
}

type Dependency struct {
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func New(d *Dependency) *Mobile {
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on attendance/mobile module")
	}
	if d.UserService == nil {
		log.Fatal("[x] UserService required on attendance/mobile module")
	}

	dbImpl := db.New(&db.Dependency{
		PgResolver: d.PgResolver,
	})

	return &Mobile{
		db:          dbImpl,
		userService: d.UserService,
	}
}

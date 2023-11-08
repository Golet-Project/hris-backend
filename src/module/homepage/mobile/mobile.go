package mobile

import (
	"hris/module/homepage/mobile/db"
	"hris/module/shared/postgres"
	"log"

	userService "hris/module/user/service"
)

type Mobile struct {
	db *db.Db

	userService *userService.Service
}

type Dependency struct {
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func New(d *Dependency) *Mobile {
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on homepage/mobile package")
	}
	if d.UserService == nil {
		log.Fatal("[x] User service required on homepage/mobile package")
	}

	dbImpl := db.New(&db.Dependency{
		PgResolver: d.PgResolver,
	})

	return &Mobile{
		db: dbImpl,

		userService: d.UserService,
	}
}

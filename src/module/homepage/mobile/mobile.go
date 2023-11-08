package mobile

import (
	"hris/module/shared/postgres"
	"log"

	userService "hris/module/user/service"
)

type Mobile struct {
	pgResolver *postgres.Resolver
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

	return &Mobile{
		pgResolver: d.PgResolver,
	}
}

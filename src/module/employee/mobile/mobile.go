package mobile

import (
	"hris/module/employee/mobile/db"
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	userService "hris/module/user/service"
)

type Mobile struct {
	db         *db.Db
	pgResolver *postgres.Resolver

	// other service
	userService *userService.Service
}

type Dependency struct {
	MasterConn  *pgxpool.Pool
	PgResolver  *postgres.Resolver
	UserService *userService.Service
}

func New(d Dependency) *Mobile {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on employee/mobile module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on employee/mobile module")
	}
	if d.UserService == nil {
		log.Fatal("[x] UserService module required on employee/mobile module")
	}

	dbImpl := db.New(db.Dependency{
		MasterConn: d.MasterConn,
		PgResolver: d.PgResolver,
	})

	return &Mobile{
		db:         dbImpl,
		pgResolver: d.PgResolver,

		userService: d.UserService,
	}
}

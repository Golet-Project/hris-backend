package mobile

import (
	"hroost/module/attendance/mobile/db"
	"hroost/module/shared/postgres"
	"log"

	userService "hroost/module/user/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Mobile struct {
	db *db.Db

	// other service
	userService *userService.Service
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func New(d *Dependency) *Mobile {
	if d.MasterConn == nil {
		log.Fatal("[x] Master connection required on attendance/mobile module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on attendance/mobile module")
	}
	if d.UserService == nil {
		log.Fatal("[x] UserService required on attendance/mobile module")
	}

	dbImpl := db.New(&db.Dependency{
		MasterConn: d.MasterConn,
		PgResolver: d.PgResolver,
	})

	return &Mobile{
		db:          dbImpl,
		userService: d.UserService,
	}
}

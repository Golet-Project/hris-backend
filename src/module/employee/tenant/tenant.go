package tenant

import (
	"hroost/module/employee/tenant/db"
	"hroost/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	db *db.Db
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	PgResolver *postgres.Resolver
}

func New(d *Dependency) *Tenant {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on employee/tenant module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Postgres resolver required on employee/tenant module")
	}

	dbImpl := db.New(&db.Dependency{
		MasterConn: d.MasterConn,
		PgResolver: d.PgResolver,
	})

	return &Tenant{
		db: dbImpl,
	}
}

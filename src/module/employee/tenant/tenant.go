package tenant

import (
	"hris/module/employee/tenant/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	db *db.Db
}

type Dependency struct {
	MasterConn *pgxpool.Pool
}

func New(d *Dependency) *Tenant {
	if d.MasterConn == nil {
		panic("[x] Master database connection required on employee/tenant module")
	}

	dbImpl := db.New(&db.Dependency{
		MasterConn: d.MasterConn,
	})

	return &Tenant{
		db: dbImpl,
	}
}

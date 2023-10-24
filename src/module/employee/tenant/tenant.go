package tenant

import (
	"hris/module/employee/tenant/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	db *db.Db
}

type Dependency struct {
	Pg *pgxpool.Pool
}

func New(d *Dependency) *Tenant {
	return &Tenant{
		db: &db.Db{
			Pg: d.Pg,
		},
	}
}

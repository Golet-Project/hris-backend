package db

import "github.com/jackc/pgx/v5/pgxpool"

type Db struct {
	masterConn *pgxpool.Pool
}

type Dependency struct {
	MasterConn *pgxpool.Pool
}

func New(d *Dependency) *Db {
	if d.MasterConn == nil {
		panic("[x] Master database connection required on employee/tenant/db module")
	}

	return &Db{
		masterConn: d.MasterConn,
	}
}

package db

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	masterConn *pgxpool.Pool
}

type Dependency struct {
	MasterConn *pgxpool.Pool
}

func New(d *Dependency) *Db {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on tenant/central/db module")
	}

	return &Db{
		masterConn: d.MasterConn,
	}
}

package db

import (
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	masterConn *pgxpool.Pool
	pgResolver *postgres.Resolver
}

type Dependency struct {
	MasterConn *pgxpool.Pool	

	PgResolver *postgres.Resolver
}

func New(d Dependency) *Db {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on employee/mobile/db module")
	}

	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on employee/mobile/db module")
	}

	return &Db {
		masterConn: d.MasterConn,
		pgResolver: d.PgResolver,
	}
}
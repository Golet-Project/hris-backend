package db

import (
	"hris/module/shared/postgres"
	"log"
)

type Db struct {
	pgResolver *postgres.Resolver
}

type Dependency struct {
	PgResolver *postgres.Resolver
}

func New(d *Dependency) *Db {
	if d.PgResolver == nil {
		log.Fatal("[x] postgres resolver is required on tenant/db module")
	}

	return &Db{
		pgResolver: d.PgResolver,
	}
}

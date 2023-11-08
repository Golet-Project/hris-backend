package db

import "hris/module/shared/postgres"

type Db struct {
	pgResolver *postgres.Resolver
}

type Dependency struct {
	PgResolver *postgres.Resolver
}

func New(d *Dependency) *Db {
	if d.PgResolver == nil {
		panic("[x] Database resolver required on homepage/mobile/db package")
	}

	return &Db{
		pgResolver: d.PgResolver,
	}
}
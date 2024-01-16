package db

import (
	"hroost/infrastructure/store/postgres"
	"log"

	"github.com/redis/go-redis/v9"
)

type Db struct {
	pgResolver *postgres.Resolver
	redis      *redis.Client
}

type Dependency struct {
	PgResolver *postgres.Resolver
	Redis      *redis.Client
}

func New(d *Dependency) *Db {
	if d.PgResolver == nil {
		log.Fatal("[x] Master database connection required on auth/central/db module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth/central/db module")
	}

	return &Db{
		pgResolver: d.PgResolver,
		redis:      d.Redis,
	}
}

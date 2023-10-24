package db

import (
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Db struct {
	masterConn *pgxpool.Pool
	pgResolver *postgres.Resolver
	redis      *redis.Client
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	PgResolver *postgres.Resolver
	Redis      *redis.Client
}

func New(d *Dependency) *Db {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth module")
	}

	return &Db{
		masterConn: d.MasterConn,
		pgResolver: d.PgResolver,
		redis:      d.Redis,
	}
}

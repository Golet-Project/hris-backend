package db

import (
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
)

type Db struct {
	masterConn *pgxpool.Pool
	pgResolver *postgres.Resolver
	redis      *redisClient.Client
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	PgResolver *postgres.Resolver
	Redis      *redisClient.Client
}

func New(d *Dependency) *Db {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth/tenant/db module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth/tenant/db module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth/tenant/db module")
	}

	return &Db{
		masterConn: d.MasterConn,
		pgResolver: d.PgResolver,
		redis:      d.Redis,
	}
}

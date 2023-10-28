package db

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Db struct {
	masterConn *pgxpool.Pool
	redis      *redis.Client
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	Redis      *redis.Client
}

func New(d *Dependency) *Db {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth/internal/db module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth/internal/db module")
	}

	return &Db{
		masterConn: d.MasterConn,
		redis:      d.Redis,
	}
}

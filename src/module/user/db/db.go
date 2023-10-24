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
		log.Fatal("[x] Database connection required on user module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on user module")
	}

	return &Db{
		masterConn: d.MasterConn,
		redis:      d.Redis,
	}
}

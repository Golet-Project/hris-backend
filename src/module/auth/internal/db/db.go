package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Db struct {
	Pg    *pgxpool.Pool
	Redis *redis.Client
}

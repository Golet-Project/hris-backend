package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	 DB *pgxpool.Pool
	 Redis *redis.Client
}
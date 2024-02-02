package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"

	redisClient "github.com/redis/go-redis/v9"
)

type IDbStore interface {
	GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error)
}

type Db struct {
	pgResolver *postgres.Resolver
	redis      *redisClient.Client
}

type Config struct {
	PgResolver *postgres.Resolver
	Redis      *redisClient.Client
}

func New(cfg *Config) (*Db, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.PgResolver == nil {
		return nil, fmt.Errorf("pgResolver required")
	}
	if cfg.Redis == nil {
		return nil, fmt.Errorf("redis required")
	}

	return &Db{
		pgResolver: cfg.PgResolver,
		redis:      cfg.Redis,
	}, nil
}

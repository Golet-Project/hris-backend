package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"

	"github.com/redis/go-redis/v9"
)

type IDbStore interface {
	ChangePassword(ctx context.Context, in ChangePasswordIn) (rowsAffected int64, err error)

	GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error)
}

type Db struct {
	pgResolver *postgres.Resolver
	redis      *redis.Client
}

type Config struct {
	PgResolver *postgres.Resolver
	Redis      *redis.Client
}

func New(cfg *Config) (*Db, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db required")
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

package db

import (
	"fmt"
	"hroost/infrastructure/store/postgres"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	PgResolver *postgres.Resolver
	Redis      *redis.Client
}

type Db struct {
	pgResolver *postgres.Resolver
	redis      *redis.Client
}

func New(cfg *Config) (*Db, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.PgResolver == nil {
		return nil, fmt.Errorf("pgResolver required")
	}
	if cfg.Redis == nil {
		return nil, fmt.Errorf("redis required at tenant.employee.db package")
	}

	return &Db{
		pgResolver: cfg.PgResolver,
		redis:      cfg.Redis,
	}, nil
}

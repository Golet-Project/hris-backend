package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"
)

type IDbStore interface {
	GetEmployeeDetail(ctx context.Context, domain, uid string) (out GetEmployeeDetailOut, err error)
}

type Db struct {
	pgResolver *postgres.Resolver
}

type Config struct {
	PgResolver *postgres.Resolver
}

func New(cfg *Config) (*Db, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.PgResolver == nil {
		return nil, fmt.Errorf("pgResolver required")
	}
	return &Db{
		pgResolver: cfg.PgResolver,
	}, nil
}

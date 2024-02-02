package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"
)

type IDbStore interface {
	CreateEmployee(ctx context.Context, data CreateEmployeeIn) (err error)

	FindAllEmployee(ctx context.Context, domain string) (out []FindAllEmployeeOut, err error)
}

type Config struct {
	PgResolver *postgres.Resolver
}

type Db struct {
	pgResolver *postgres.Resolver
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

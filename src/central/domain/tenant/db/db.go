package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"
)

type IDbStore interface {
	CreateTenant(ctx context.Context, in CreateTenantIn) (out CreateTenantOut, err error)

	CountTenantByDomain(ctx context.Context, domain string) (out CountTenantByDomainOut, err error)
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

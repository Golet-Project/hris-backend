package village

import (
	"fmt"
	"hroost/infrastructure/store/postgres"
)

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
		return nil, fmt.Errorf("pgResovler required")
	}

	return &Db{
		pgResolver: cfg.PgResolver,
	}, nil
}

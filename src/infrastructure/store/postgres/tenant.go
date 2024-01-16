package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantDbConfig struct {
	Domain   string
	User     string
	Password string
	Host     string
	Port     string
}

type tenantDb struct {
	connString string
}

func NewTenantDb(cfg *TenantDbConfig) (*tenantDb, error) {
	if cfg == nil {
		return nil, fmt.Errorf("tenantDbConfig can not be empty")
	}
	if cfg.Domain == "" {
		return nil, fmt.Errorf("Domain name can not be empty")
	}

	dbName := "tenant_" + cfg.Domain
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		dbName,
	)
	return &tenantDb{
		connString: connString,
	}, nil
}

func (t *tenantDb) config() (*pgxpool.Config, error) {
	connConfig, err := pgxpool.ParseConfig(t.connString)
	if err != nil {
		return nil, err
	}

	connConfig.MinConns = 3
	connConfig.MaxConns = 5
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	return connConfig, nil
}

func (t *tenantDb) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connConfig, err := t.config()
	if err != nil {
		return nil, err
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	// ping
	err = pgPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pgPool, nil
}

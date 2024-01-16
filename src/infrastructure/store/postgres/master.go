package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MasterDbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
}

type masterDb struct {
	connString string
}

func NewMasterDb(cfg *MasterDbConfig) (*masterDb, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can not be empty")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
	)

	return &masterDb{
		connString: connString,
	}, nil
}

func (m *masterDb) config() (*pgxpool.Config, error) {
	connConfig, err := pgxpool.ParseConfig(m.connString)
	if err != nil {
		return nil, err
	}

	connConfig.MinConns = 3
	connConfig.MaxConns = 5
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	return connConfig, nil
}

func (m *masterDb) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connConfig, err := m.config()
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

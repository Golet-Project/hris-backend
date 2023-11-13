package postgres

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MasterConnConfig() (*pgxpool.Config, error) {
	var connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_MASTER_DB"),
	)

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connConfig.MinConns = 3
	connConfig.MaxConns = 5
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	return connConfig, nil
}

func TenantConnConfig(dbName string) (*pgxpool.Config, error) {
	var connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		dbName,
	)

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connConfig.MinConns = 3
	connConfig.MaxConns = 5
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	return connConfig, nil
}

func WorkerConnConfig() (*pgxpool.Config, error) {
	var connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_MASTER_DB"),
	)

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = 3
	connConfig.MinConns = 1
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	return connConfig, nil
}

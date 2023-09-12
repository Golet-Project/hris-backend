package tenant

import (
	"context"
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	workerDBConn *pgxpool.Pool
	PgResolver   *postgres.Resolver
}

func NewMigration(ctx context.Context, workerDBConn *pgxpool.Pool, pgResolver *postgres.Resolver) (*Migration, error) {
	return &Migration{
		workerDBConn: workerDBConn,
		PgResolver:   pgResolver,
	}, nil
}

type Tenant struct {
	Domain string
}

func (m *Migration) Run(ctx context.Context, tenant Tenant) error {
	// create database
	dbName, err := m.CreateDatabase(ctx, tenant.Domain)
	if err != nil {
		return err
	}

	connConfig, err := postgres.TenantConnConfig(dbName)
	if err != nil {
		log.Println("[x] Failed to make tenant database connection")
		return err
	}

	pgPool, err := pgxpool.NewWithConfig(context.TODO(), connConfig)
	if err != nil {
		return err
	}

	if err := pgPool.Ping(context.TODO()); err != nil {
		log.Println("[x] Failed to make tenant database connection")
		return err
	}

	// store to resolver
	m.PgResolver.Register(postgres.Database{
		DomainName: postgres.Domain(tenant.Domain),
		Pool:       pgPool,
	})

	tx, err := pgPool.Begin(context.TODO())

	if err := m.CreateExtensionUuid(context.TODO(), tx); err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return e
		}

		return err
	}

	// migrate employee table
	if err := m.CreateTableEmployee(context.TODO(), tx); err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return e
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// TODO: delete this
	log.Print("=====\n\n")
	log.Println("current connected tenant is:")
	m.PgResolver.GetAllTenantConn().Range(func(key, val interface{}) bool {
		log.Printf("%s", key)
		return true
	})
	log.Print("=====\n\n")

	return nil
}

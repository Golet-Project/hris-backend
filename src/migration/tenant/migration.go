package tenant

import (
	"context"
	"fmt"
	"hroost/module/shared/postgres"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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

	migrationDatabaseURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		dbName,
	)
	goMigrate, err := migrate.New("file://migration/tenant/postgres", migrationDatabaseURL)
	if err != nil {
		return err
	}
	defer goMigrate.Close()

	log.Println("running database migration...")
	err = goMigrate.Up()
	if err != nil {
		return err
	}

	tenantConnConfig, err := postgres.TenantConnConfig(dbName)
	if err != nil {
		log.Println("[x] Failed to make tenant database connection")
		return err
	}

	pgPool, err := pgxpool.NewWithConfig(ctx, tenantConnConfig)
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

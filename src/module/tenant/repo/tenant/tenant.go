package tenant

import (
	"hris/module/shared/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
	PgResolver *postgres.Resolver
}
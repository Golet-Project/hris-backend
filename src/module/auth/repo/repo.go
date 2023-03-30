package repo

import "github.com/jackc/pgx/v5/pgxpool"

type AuthRepo struct {
	DB *pgxpool.Pool
}

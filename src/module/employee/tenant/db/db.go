package db

import "github.com/jackc/pgx/v5/pgxpool"

type Db struct {
	Pg 	*pgxpool.Pool
}
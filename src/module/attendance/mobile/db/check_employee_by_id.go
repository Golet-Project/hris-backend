package db

import (
	"context"
	"hroost/module/shared/postgres"

	"github.com/jackc/pgx/v5"
)

func (d *Db) CheckEmployeeById(ctx context.Context, domain string, uid string) (exist bool, err error) {
	if domain == "" || uid == "" {
		return false, pgx.ErrNoRows
	}

	var count int64
	var sql = `SELECT COUNT(id) FROM employee WHERE uid = $1`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	err = conn.QueryRow(ctx, sql, uid).Scan(&count)
	if err != nil {
		return
	}

	return true, nil
}

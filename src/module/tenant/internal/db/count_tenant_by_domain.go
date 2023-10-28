package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type CountTenantByDomainOut struct {
	Count int64
}

func (d *Db) CountTenantByDomain(ctx context.Context, domain string) (out CountTenantByDomainOut, err error) {
	sql := `SELECT COUNT(id) FROM tenant WHERE domain = @domain`

	err = d.masterConn.QueryRow(ctx, sql, pgx.NamedArgs{
		"domain": domain,
	}).Scan(&out.Count)

	return
}

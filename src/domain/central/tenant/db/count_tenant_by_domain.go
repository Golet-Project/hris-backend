package db

import (
	"context"
	"hroost/infrastructure/store/postgres"

	"github.com/jackc/pgx/v5"
)

type CountTenantByDomainOut struct {
	Count int64
}

func (d *Db) CountTenantByDomain(ctx context.Context, domain string) (out CountTenantByDomainOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, err
	}

	sql := `SELECT COUNT(id) FROM tenant WHERE domain = @domain`

	err = masterConn.QueryRow(ctx, sql, pgx.NamedArgs{
		"domain": domain,
	}).Scan(&out.Count)

	return
}

package tenant

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type CountTenantByDomainOut struct {
	Count int64
}

func (r *Repository) CountTenantByDomain(ctx context.Context, domain string) (out CountTenantByDomainOut, err error) {
	sql := `SELECT COUNT(id) FROM tenant WHERE domain = @domain`

	err = r.DB.QueryRow(ctx, sql, pgx.NamedArgs{
		"domain": domain,
	}).Scan(&out.Count)

	return
}

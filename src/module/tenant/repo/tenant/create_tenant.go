package tenant

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Internal_CreateTenantIn struct {
	Name   string
	Domain string
}

type Interal_CreateTenantOut struct {
	UID    string
	Name   string
	Domain string
}

func (r *Repository) Internal_CreateTenant(ctx context.Context, in Internal_CreateTenantIn) (out Interal_CreateTenantOut, err error) {
	sql := `
	INSERT INTO tenant (
		name, domain
	) VALUES (
		@name, @domain
	) RETURNING uid, name, domain;`

	err = r.DB.QueryRow(ctx, sql, pgx.NamedArgs{
		"name":   in.Name,
		"domain": in.Domain,
	}).Scan(&out.UID, &out.Name, &out.Domain)

	return
}

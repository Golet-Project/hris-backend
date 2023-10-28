package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type CreateTenantIn struct {
	Name   string
	Domain string
}

type CreateTenantOut struct {
	UID    string
	Name   string
	Domain string
}

func (d *Db) CreateTenant(ctx context.Context, in CreateTenantIn) (out CreateTenantOut, err error) {
	sql := `
	INSERT INTO tenant (
		name, domain
	) VALUES (
		@name, @domain
	) RETURNING uid, name, domain;`

	err = d.masterConn.QueryRow(ctx, sql, pgx.NamedArgs{
		"name":   in.Name,
		"domain": in.Domain,
	}).Scan(&out.UID, &out.Name, &out.Domain)

	return
}

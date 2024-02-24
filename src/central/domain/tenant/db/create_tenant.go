package db

import (
	"context"
	"errors"
	"hroost/central/domain/tenant/model"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) CreateTenant(ctx context.Context, in model.CreateTenantIn) (out model.CreateTenantOut, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	sql := `
	INSERT INTO tenant (
		name, domain
	) VALUES (
		@name, @domain
	) RETURNING uid, name, domain;`

	err = masterConn.QueryRow(ctx, sql, pgx.NamedArgs{
		"name":   in.Name,
		"domain": in.Domain,
	}).Scan(&out.UID, &out.Name, &out.Domain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	return
}

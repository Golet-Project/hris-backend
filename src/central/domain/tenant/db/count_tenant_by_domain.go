package db

import (
	"context"
	"errors"
	"hroost/central/domain/tenant/model"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) CountTenantByDomain(ctx context.Context, domain string) (out model.CountTenantByDomainOut, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	sql := `SELECT COUNT(id) FROM tenant WHERE domain = @domain`

	err = masterConn.QueryRow(ctx, sql, pgx.NamedArgs{
		"domain": domain,
	}).Scan(&out.Count)
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

package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
	"hroost/tenant/domain/auth/model"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out model.GetLoginCredentialOut, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	var sql = `
	SELECT
		uid, email, password, domain
	FROM
		tenant_admin
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserID, &out.Email, &out.Password, &out.Domain,
	)
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

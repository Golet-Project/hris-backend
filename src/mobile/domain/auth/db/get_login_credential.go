package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/auth/model"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out model.GetLoginCredentialOut, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	var sql = `
	SELECT
		uid, email, password, domain
	FROM
		users
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserUID, &out.Email, &out.Password, &out.Domain,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	return out, nil
}

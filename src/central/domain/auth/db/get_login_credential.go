package db

import (
	"context"
	"errors"
	"hroost/central/domain/auth/model"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetLoginCredential(ctx context.Context, email string) (model.GetLoginCredentialOut, *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	var out = model.GetLoginCredentialOut{}

	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	var sql = `
	SELECT
		uid, email, password
	FROM
		central_admin
	WHERE
	 	email = $1`

	err = masterConn.QueryRow(ctx, sql, email).Scan(&out.UserUID, &out.Email, &out.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
			}

		} else {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
			}
		}
	}

	return out, nil
}

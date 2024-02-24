package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/auth/model"
	"hroost/shared/primitive"
)

func (d *Db) ChangePassword(ctx context.Context, in model.ChangePasswordIn) (rowsAffected int64, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return rowsAffected, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	sql := `
	UPDATE
		users
	SET
		password = $1
	WHERE
		uid = $2
	`

	commandTag, err := masterConn.Exec(ctx, sql, in.Password, in.UID)
	if err != nil {
		return rowsAffected, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

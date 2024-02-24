package db

import (
	"context"
	"hroost/central/domain/auth/model"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) ChangePassword(ctx context.Context, in model.ChangePasswordIn) (int64, *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return 0, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	sql := `
	UPDATE central_admin SET
		password = @password
	WHERE
		uid = @uid`

	commandTag, err := masterConn.Exec(ctx, sql, pgx.NamedArgs{
		"password": in.Password,
		"uid":      in.UID,
	})
	if err != nil {
		return 0, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	rowsAffected := commandTag.RowsAffected()

	return rowsAffected, nil
}

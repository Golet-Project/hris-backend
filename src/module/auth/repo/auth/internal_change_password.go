package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type InternalChangePasswordIn struct {
	UID      string
	Password string
}

func (r *Repository) InternalChangePassword(ctx context.Context, in InternalChangePasswordIn) (rowsAffected int64, err error) {
	sql := `
	UPDATE internal_admin SET
		password = @password
	WHERE
		uid = @uid`

	commandTag, err := r.DB.Exec(ctx, sql, pgx.NamedArgs{
		"password": in.Password,
		"uid"	: in.UID,
	})
	if err != nil {
		return
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

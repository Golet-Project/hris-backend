package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ChangePasswordIn struct {
	UID      string
	Password string
}

func (d *Db) ChangePassword(ctx context.Context, in ChangePasswordIn) (rowsAffected int64, err error) {
	sql := `
	UPDATE central_admin SET
		password = @password
	WHERE
		uid = @uid`

	commandTag, err := d.masterConn.Exec(ctx, sql, pgx.NamedArgs{
		"password": in.Password,
		"uid":      in.UID,
	})
	if err != nil {
		return
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

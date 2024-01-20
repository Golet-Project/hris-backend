package db

import (
	"context"
	"hroost/infrastructure/store/postgres"

	"github.com/jackc/pgx/v5"
)

type ChangePasswordIn struct {
	UID      string
	Password string
}

func (d *Db) ChangePassword(ctx context.Context, in ChangePasswordIn) (rowsAffected int64, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return 0, err
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
		return
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

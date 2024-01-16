package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
)

type ChangePasswordIn struct {
	UID      string
	Password string
}

func (d *Db) ChangePassword(ctx context.Context, in ChangePasswordIn) (rowsAffected int64, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
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
		return
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

package db

import (
	"context"
)

type ChangePasswordIn struct {
	UID      string
	Password string
}

func (d *Db) ChangePassword(ctx context.Context, in ChangePasswordIn) (rowsAffected int64, err error) {
	sql := `
	UPDATE
		users
	SET
		password = $1
	WHERE
		uid = $2
	`

	commandTag, err := d.masterConn.Exec(ctx, sql, in.Password, in.UID)
	if err != nil {
		return
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

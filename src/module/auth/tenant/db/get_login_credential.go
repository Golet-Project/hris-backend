package db

import (
	"context"
)

type GetLoginCredentialOut struct {
	UserID   string
	Email    string
	Password string
}

func (d *Db) GetLoginCredential(ctx context.Context, email, domain string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password
	FROM
		tenant_admin
	WHERE
		email = $1
		AND
		domain = $2
		AND
		deleted_at IS NULL`

	err = d.masterConn.QueryRow(ctx, sql, email, domain).Scan(
		&out.UserID, &out.Email, &out.Password,
	)

	return
}

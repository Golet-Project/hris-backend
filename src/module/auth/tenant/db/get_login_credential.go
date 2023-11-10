package db

import (
	"context"
)

type GetLoginCredentialOut struct {
	UserID   string
	Email    string
	Password string
	Domain   string
}

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password, domain
	FROM
		tenant_admin
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = d.masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserID, &out.Email, &out.Password, &out.Domain,
	)

	return
}

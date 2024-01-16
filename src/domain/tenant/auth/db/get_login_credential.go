package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
)

type GetLoginCredentialOut struct {
	UserID   string
	Email    string
	Password string
	Domain   string
}

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	var sql = `
	SELECT
		uid, email, password, domain
	FROM
		tenant_admin
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserID, &out.Email, &out.Password, &out.Domain,
	)

	return
}

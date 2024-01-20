package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
)

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password primitive.String
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
		users
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserUID, &out.Email, &out.Password, &out.Domain,
	)

	return
}

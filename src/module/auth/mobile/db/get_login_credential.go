package db

import (
	"context"
	"hris/module/shared/primitive"
)

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password primitive.String
	Domain   string
}

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password, domain
	FROM
		users
	WHERE
		email = $1
		AND
		deleted_at IS NULL`

	err = d.masterConn.QueryRow(ctx, sql, email).Scan(
		&out.UserUID, &out.Email, &out.Password, &out.Domain,
	)

	return
}

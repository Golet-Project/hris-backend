package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/module/shared/primitive"
)

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password primitive.String
}

func (d *Db) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)

	var sql = `
	SELECT
		uid, email, password
	FROM
		central_admin
	WHERE
	 	email = $1`

	err = masterConn.QueryRow(ctx, sql, email).Scan(&out.UserUID, &out.Email, &out.Password)

	return
}

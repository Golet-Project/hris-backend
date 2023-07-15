package auth

import (
	"context"
	"hris/module/shared/primitive"
)

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password primitive.String
}

func (r *Repository) InternalGetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password
	FROM
		internal_admin
	WHERE
	 	email = $1`

	err = r.DB.QueryRow(ctx, sql, email).Scan(&out.UserUID, &out.Email, &out.Password)

	return
}

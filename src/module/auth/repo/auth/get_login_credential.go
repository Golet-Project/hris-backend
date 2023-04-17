package auth

import "context"

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password string
}

func (r *Repository) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password
	FROM
		users
	WHERE
	 	email = $1`

	err = r.DB.QueryRow(ctx, sql, email).Scan(&out.UserUID, &out.Email, &out.Password)

	return
}

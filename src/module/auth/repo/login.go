package repo

import "context"

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password string
}

func (a *AuthRepo) GetLoginCredential(ctx context.Context, email string) (out GetLoginCredentialOut, err error) {
	var sql = `
	SELECT
		uid, email, password
	FROM
		employee
	WHERE
	 	email = $1`

	err = a.DB.QueryRow(ctx, sql, email).Scan(&out.UserUID, &out.Email, &out.Password)

	return
}

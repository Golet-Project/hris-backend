package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

type GetEmployeeDetailOut struct {
	UID            string
	FullName       string
	Email          string
	Gender         primitive.Gender
	BirthDate      primitive.Date
	ProfilePicture primitive.String
	Address        primitive.String
	JoinDate       primitive.Date
}

func (d *Db) GetEmployeeDetail(ctx context.Context, domain, uid string) (out GetEmployeeDetailOut, err error) {
	if domain == "" || uid == "" {
		return out, pgx.ErrNoRows
	}

	var sql = `
	SELECT
		uid, full_name, email, gender, birth_date, profile_picture,
		address, join_date
	FROM
		employee
	WHERE uid = $1`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	err = conn.QueryRow(ctx, sql, uid).Scan(
		&out.UID, &out.FullName, &out.Email, &out.Gender, &out.BirthDate,
		&out.ProfilePicture, &out.Address, &out.JoinDate,
	)

	return
}

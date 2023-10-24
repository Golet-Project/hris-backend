package db

import (
	"context"
	"hris/module/shared/postgres"
	"hris/module/shared/primitive"
)

type GetEmployeeDetailOut struct {
	Email          string
	FullName       string
	Gender         primitive.Gender
	BirthDate      primitive.Date
	ProfilePicture primitive.String
	Address        primitive.String
	JoinDate       primitive.Date
}

func (d *Db) GetEmployeeDetail(ctx context.Context, domain string, uid string) (out GetEmployeeDetailOut, err error) {
	var sql = `
	SELECT
		email, full_name, gender, birth_date, profile_picture, 
		address, join_date
	FROM
		employee
		WHERE uid = $1`

	pool, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	err = pool.QueryRow(ctx, sql, uid).Scan(
		&out.Email, &out.FullName, &out.Gender, &out.BirthDate, &out.ProfilePicture,
		&out.Address, &out.JoinDate,
	)

	return
}

package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/auth/model"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetEmployeeDetail(ctx context.Context, domain string, uid string) (out model.GetEmployeeDetailOut, repoError *primitive.RepoError) {
	var sql = `
	SELECT
		email, full_name, gender, birth_date, profile_picture, 
		address, join_date
	FROM
		employee
		WHERE uid = $1`

	pool, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	err = pool.QueryRow(ctx, sql, uid).Scan(
		&out.Email, &out.FullName, &out.Gender, &out.BirthDate, &out.ProfilePicture,
		&out.Address, &out.JoinDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
			}
		}

		return out, &primitive.RepoError{Issue: primitive.RepoErrorCodeServerError}
	}

	return out, nil
}

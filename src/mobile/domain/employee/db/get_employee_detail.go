package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/employee/model"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetEmployeeDetail(ctx context.Context, domain, uid string) (out model.GetEmployeeDetailOut, repoError *primitive.RepoError) {
	if domain == "" || uid == "" {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   pgx.ErrNoRows,
		}
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
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	err = conn.QueryRow(ctx, sql, uid).Scan(
		&out.UID, &out.FullName, &out.Email, &out.Gender, &out.BirthDate,
		&out.ProfilePicture, &out.Address, &out.JoinDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	return
}

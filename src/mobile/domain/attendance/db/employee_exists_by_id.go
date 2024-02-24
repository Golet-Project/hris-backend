package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) EmployeeExistsById(ctx context.Context, domain string, uid string) (exist bool, repoError *primitive.RepoError) {
	if domain == "" || uid == "" {
		return false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   pgx.ErrNoRows,
		}
	}

	var count int64
	var sql = `SELECT COUNT(id) FROM employee WHERE uid = $1`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	err = conn.QueryRow(ctx, sql, uid).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   pgx.ErrNoRows,
			}
		}

		return false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	return true, nil
}

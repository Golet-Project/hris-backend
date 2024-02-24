package db

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
)

func (d *Db) Checkout(ctx context.Context, domain, uid string) (rowsAffected int64, repoError *primitive.RepoError) {
	var sql = `
	UPDATE
		attendance
	SET
		checkout_time = now()
	WHERE
		employee_uid = $1
		AND
		checkout_time IS NULL`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return 0, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	commandTag, err := conn.Exec(ctx, sql, uid)
	if err != nil {
		return 0, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	rowsAffected = commandTag.RowsAffected()

	return
}

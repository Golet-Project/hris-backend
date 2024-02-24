package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
	"hroost/tenant/domain/attendance/model"

	"github.com/jackc/pgx/v5"
)

func (d *Db) FindAllAttendance(ctx context.Context, domain string) (out []model.FindAllAttendanceOut, repoError *primitive.RepoError) {
	if domain == "" {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   pgx.ErrNoRows,
		}
	}

	var sql = `
	SELECT
		a.uid, e.full_name, a.checkin_time, a.checkout_time, a.approved_at,
		a.approved_by
	FROM
		attendance AS a
		INNER JOIN employee AS e ON e.uid = a.employee_uid`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	rows, err := conn.Query(ctx, sql)
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
	defer rows.Close()

	for rows.Next() {
		var row model.FindAllAttendanceOut

		err = rows.Scan(
			&row.UID, &row.FullName, &row.CheckinTime, &row.CheckoutTime, &row.ApprovedAt,
			&row.ApprovedBy,
		)
		if err != nil {
			return out, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
				Err:   err,
			}

		}

		out = append(out, row)
	}

	return
}

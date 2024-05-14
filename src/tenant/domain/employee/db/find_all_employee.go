package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
	"hroost/tenant/domain/employee/model"

	"github.com/jackc/pgx/v5"
)

func (d *Db) FindAllEmployee(ctx context.Context, domain string) (out []model.FindAllEmployeeOut, repoError *primitive.RepoError) {
	if domain == "" {
		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   pgx.ErrNoRows,
		}
	}

	sql := `
	SELECT
		e.uid, e.email, e.full_name, e.birth_date, e.gender, e.employee_status,
		e.join_date, e.end_date
	FROM
		employee AS e`

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
				Err:   pgx.ErrNoRows,
			}
		}

		return out, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}
	defer rows.Close()

	for rows.Next() {
		var row model.FindAllEmployeeOut

		err = rows.Scan(&row.Id, &row.Email, &row.FullName, &row.BirthDate, &row.Gender, &row.EmployeeStatus,
			&row.JoinDate, &row.EndDate,
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

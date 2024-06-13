package db

import (
	"context"
	"errors"
	"fmt"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
	"hroost/tenant/domain/employee/model"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetById(ctx context.Context, domain, id string) (out model.GetByIdOut, repoErr *primitive.RepoError) {
	if id == "" {
		repoErr = &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("missing employee id"),
		}

		return
	}

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		repoErr = &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}

		return
	}

	sql := `
  SELECT 
    e.uid, e.full_name, e.email, e.gender, e.employee_status, e.birth_date,
    e.join_date, e.address, e.province_id, p.name AS province_name, e.regency_id, r.name AS regency_name,
    e.district_id, d.name AS district_name, e.village_id, v.name AS village_name
  FROM
    employee AS e
  INNER JOIN province AS p ON p.id = e.province_id
  INNER JOIN regency AS r ON r.id = e.regency_id
  INNER JOIN district AS d ON d.id = e.district_id
  INNER JOIN village AS v ON v.id = e.village_id
  WHERE
    e.uid = $1`

	err = conn.QueryRow(ctx, sql, id).Scan(
		&out.Id, &out.FullName, &out.Email, &out.Gender, &out.EmployeeStatus, &out.BirthDate, &out.JoinDate,
		&out.Address.Detail, &out.Address.ProvinceId, &out.Address.ProvinceName, &out.Address.RegencyId, &out.Address.RegencyName,
		&out.Address.DistrictId, &out.Address.DistrictName, &out.Address.VillageId, &out.Address.VillageName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repoErr = &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}

			return
		}

		repoErr = &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}

		return
	}

	return
}

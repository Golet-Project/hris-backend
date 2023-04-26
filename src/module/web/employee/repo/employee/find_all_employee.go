package employee

import (
	"context"
	"hris/module/shared/primitive"
	"time"
)

type FindAllEmployeeOut struct {
	UID          string
	Email        string
	FullName     string
	ProvinceName string
	RegencyName  string
	DistrictName string
	VillageName  string
	RegisteredAt time.Time
}

func (r Repository) FindAllEmployee(ctx context.Context) ([]FindAllEmployeeOut, error) {
	var sql = `
	SELECT
		u.uid, u.email, u.full_name,
		p.name AS province_name,
		r.name AS regency_name,
		d.name AS district_name,
		v.name AS village_name,
		u.created_at AS registered_at
	FROM
		users AS u
		INNER JOIN province AS p ON p.id = u.province_id
		INNER JOIN regency AS r ON r.id = u.regency_id
		INNER JOIN district AS d ON d.id = u.district_id
		INNER JOIN village AS v ON v.id = u.village_id
	WHERE
		u.type = $1`

	rows, err := r.DB.Query(ctx, sql, primitive.UserTypeEmployee)
	if err != nil {
		return []FindAllEmployeeOut{}, err
	}

	defer rows.Close()

	var employee = []FindAllEmployeeOut{}
	for rows.Next() {
		var e FindAllEmployeeOut
		err := rows.Scan(
			&e.UID, &e.Email, &e.FullName,
			&e.ProvinceName, &e.RegencyName, &e.DistrictName,
			&e.VillageName, &e.RegisteredAt,
		)
		if err != nil {
			return []FindAllEmployeeOut{}, err
		}

		employee = append(employee, e)
	}

	return employee, nil
}

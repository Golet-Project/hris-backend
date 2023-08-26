package province

import (
	"context"
)

type FindAllProvinceOut struct {
	ID   string
	Name string
}

func (r *Repository) FindAllProvince(ctx context.Context) (out []FindAllProvinceOut, err error) {
	var sql = `
	SELECT
		id, name
	FROM
		province`

	rows, err := r.DB.Query(ctx, sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row FindAllProvinceOut

		err = rows.Scan(&row.ID, &row.Name)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

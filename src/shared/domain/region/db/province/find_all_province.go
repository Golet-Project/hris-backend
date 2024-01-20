package province

import (
	"context"
	"hroost/infrastructure/store/postgres"
)

type FindAllProvinceOut struct {
	ID   string
	Name string
}

func (r *Db) FindAllProvince(ctx context.Context) (out []FindAllProvinceOut, err error) {
	masterConn, err := r.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return nil, err
	}

	var sql = `
	SELECT
		id, name
	FROM
		province`

	rows, err := masterConn.Query(ctx, sql)
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

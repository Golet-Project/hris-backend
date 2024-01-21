package regency

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/entities"
)

type FindAllByProvinceOut struct {
	entities.Regency
}

func (d Db) FindAllByProvinceId(ctx context.Context, provinceId string) (out []FindAllByProvinceOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	var sql = `SELECT id, province_id, name FROM regency WHERE province_id = $1`

	rows, err := masterConn.Query(ctx, sql, provinceId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row FindAllByProvinceOut

		err = rows.Scan(&row.Id, &row.ProvinceId, &row.Name)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

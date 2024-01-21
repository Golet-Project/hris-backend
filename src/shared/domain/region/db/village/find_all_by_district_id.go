package village

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/entities"
)

type FindAllByDistrictId struct {
	entities.Village
}

func (d Db) FindAllByDistrictId(ctx context.Context, provinceId string) (out []FindAllByDistrictId, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	var sql = `SELECT id, district_id, name FROM village WHERE district_id = $1`

	rows, err := masterConn.Query(ctx, sql, provinceId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row FindAllByDistrictId

		err = rows.Scan(&row.Id, &row.DistrictId, &row.Name)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

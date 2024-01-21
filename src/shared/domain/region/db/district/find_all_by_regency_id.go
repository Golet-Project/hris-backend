package district

import (
	"context"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/entities"
)

type FindAllByRegencyIdOut struct {
	entities.District
}

func (d Db) FindAllByRegencyId(ctx context.Context, regencyId string) (out []FindAllByRegencyIdOut, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	var sql = `SELECT id, regency_id, name FROM district WHERE regency_id = $1`

	rows, err := masterConn.Query(ctx, sql, regencyId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var row FindAllByRegencyIdOut

		err = rows.Scan(&row.Id, &row.RegencyId, &row.Name)
		if err != nil {
			return
		}

		out = append(out, row)
	}

	return
}

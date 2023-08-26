package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateTableAndSeedDistrict(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS district (
		id VARCHAR(8) NOT NULL,
		regency_id VARCHAR(5) NOT NULL,
		name VARCHAR(255) NOT NULL,

		PRIMARY KEY (id)
	);

	CREATE INDEX IF NOT EXISTS district_regency_id_idx ON district(regency_id);
	
	INSERT INTO district (id, regency_id, name)
	SELECT
    kode AS id,
    LEFT(kode, 5) AS regency_id,
    nama
	FROM wilayah
	WHERE LENGTH(kode) = 8;`

	log.Println("CREATE and SEED district TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

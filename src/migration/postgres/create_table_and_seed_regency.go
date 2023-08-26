package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateTableAndSeedRegency(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS regency (
		id VARCHAR(5) NOT NULL,
		province_id VARCHAR(2) NOT NULL,
		name VARCHAR(255) NOT NULL,

		PRIMARY KEY (id)
	);

	CREATE INDEX IF NOT EXISTS regency_province_id_idx ON regency(province_id);

	INSERT INTO regency (id, province_id, name)
	SELECT
    kode AS id,
    LEFT(kode, 2) AS province_id,
    nama
	FROM wilayah
	WHERE LENGTH(kode) = 5;`

	log.Println("CREATE and SEED regency TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

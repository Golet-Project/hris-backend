package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateTableAndSeedVillage(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS village (
		id VARCHAR(13) NOT NULL,
		district_id VARCHAR(8) NOT NULL,
		name VARCHAR(255) NOT NULL,

		PRIMARY KEY (id)
	);

	CREATE INDEX IF NOT EXISTS village_district_id_idx ON village(district_id);
	
	INSERT INTO village (id, district_id, name)
	SELECT
    kode AS id,
    LEFT(kode, 8) AS district_id,
    nama
	FROM wilayah
	WHERE LENGTH(kode) = 13;`

	log.Println("CREATE and SEED village TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

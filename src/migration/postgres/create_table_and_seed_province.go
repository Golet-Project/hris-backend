package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateTableAndSeedProvince(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS province (
    id VARCHAR(2) NOT NULL,
    name varchar(255) not null,

    PRIMARY KEY (id)
	);

	INSERT INTO province (id, name)
	SELECT
    kode AS id,
    nama
	FROM wilayah
	WHERE LENGTH(kode) = 2;`

	log.Println("CREATE and SEED province TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

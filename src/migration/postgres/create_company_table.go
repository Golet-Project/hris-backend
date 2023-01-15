package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateCompanyTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS company (
		pk_id		BIGSERIAL NOT NULL,
		id			CHAR(36) NOT NULL,
		detail		TEXT NOT NULL,
		latitude	FLOAT8 NOT NULL,
		longitude 	FLOAT8 NOT NULL,
		address		VARCHAR(500) NOT NULL,
		province_id	INTEGER	NOT NULL,
		regency_id	INTEGER NOT NULL,
		district_id	INTEGER NOT NULL,
		village_id	INTEGER NOT NULL,
		created_at	BIGINT NOT NULL,
		updated_at	BIGINT NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id)
	)`

	log.Println("CREATING company TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

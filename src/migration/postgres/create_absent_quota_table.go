package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAbsentQuotaTable(ctx context.Context) error {
		var sql = `
	CREATE TABLE IF NOT EXISTS absent_quota (
		pk_id			BIGSERIAL NOT NULL,
		id				CHAR(36) NOT NULL,
		user_company_id	CHAR(36) NOT NULL,
		quota			INTEGER NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id)
	)`

	log.Println("CREATING absent_quota TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

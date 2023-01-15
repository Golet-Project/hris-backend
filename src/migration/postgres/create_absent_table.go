package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAbsentTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS absent (
		pk_id				BIGSERIAL NOT NULL,
		id					CHAR(36) NOT NULL,
		user_id				CHAR(36) NOT NULL,
		company_id			CHAR(36) NOT NULL,
		absent_date			DATE NOT NULL,
		reason				VARCHAR(500) NOT NULL,
		type				absent_type NOT NULL,
		created_at			BIGINT NOT NULL,
		created_by_user_id	BIGINT NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (user_id, company_id, absent_date)
	)`

	log.Println("CREATING absent TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

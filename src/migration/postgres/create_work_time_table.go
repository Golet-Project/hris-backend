package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateWorkTimeTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS work_time (
		pk_id		BIGSERIAL NOT NULL,
		id 			CHAR(36) NOT NULL,
		company_id	CHAR(36) NOT NULL,
		day			SMALLINT NOT NULL,
		shift		SMALLINT NOT NULL,
		time_start	SMALLINT NOT NULL,
		time_end	SMALLINT NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id)
	)`

	log.Println("CREATING work_time TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateRoleTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS role (
		pk_id			BIGSERIAL NOT NULL,
		id				CHAR(36) NOT NULL,
		company_id		CHAR(36) NOT NULL,
		name			VARCHAR(255) NOT NULL,
		created_at		BIGINT NOT NULL,
		created_by_user_id	CHAR(36) NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (company_id, name),
		UNIQUE (company_id)
	)`

	log.Println("CREATING role TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

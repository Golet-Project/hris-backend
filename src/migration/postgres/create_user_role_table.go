package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateUserRoleTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS user_role (
		pk_id			BIGSERIAL NOT NULL,
		id				CHAR(36) NOT NULL,
		user_company_id	CHAR(36) NOT NULL,
		role_id			CHAR(36) NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (user_company_id, role_id)
	)`

	log.Println("CREATING user_role TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

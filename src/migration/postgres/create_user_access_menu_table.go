package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateUserAccessMenuTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS user_access_menu (
		id					BIGSERIAL NOT NULL,
		user_company_id		CHAR(36) NOT NULL,
		access_menu_id		INTEGER NOT NULL,
		created_at			BIGINT NOT NULL,
		created_by_user_id	CHAR(36) NOT NULL,

		PRIMARY KEY (id),
		UNIQUE (user_company_id, access_menu_id)
	)`

	log.Println("CREATING user_access_menu TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

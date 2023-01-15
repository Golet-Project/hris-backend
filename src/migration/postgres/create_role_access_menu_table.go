package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateRoleAccessMenuTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS role_access_menu (
		id					BIGSERIAL NOT NULL,
		role_id				CHAR(36) NOT NULL,
		access_menu_id		INTEGER NOT NULL,
		created_at			BIGINT NOT NULL,
		created_by_user_id	BIGINT NOT NULL,

		PRIMARY KEY (id),
		UNIQUE (role_id, access_menu_id)
	)`

	log.Println("CREATING role_access_menu TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

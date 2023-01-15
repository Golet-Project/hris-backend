package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateUserCompanyTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS user_company (
		pk_id		BIGSERIAL NOT NULL,
		id			CHAR(36) NOT NULL,
		user_id		CHAR(36) NOT NULL,
		company_id	CHAR(36) NOT NULL,
		joined_date	DATE NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (user_id, company_id)
	)`

	log.Println("CREATING user_company TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

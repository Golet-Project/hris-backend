package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateUserCompanyWorkTimeTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS user_company_work_time (
		pk_id			SERIAL NOT NULL,
		id				CHAR(36) NOT NULL,
		user_company_id	CHAR(36) NOT NULL,
		work_time_id	CHAR(36) NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (user_company_id, work_time_id)
	)`

	log.Println("CREATING user_company_work_time TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

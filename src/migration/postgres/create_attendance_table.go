package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAttendanceTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS attendance (
		pk_id				BIGSERIAL NOT NULL,
		id					CHAR(36) NOT NULL,
		user_id				CHAR(36) NOT NULL,
		company_id			CHAR(36) NOT NULL,
		working_date		CHAR(36) NOT NULL,
		work_time_id		CHAR(36) NOT NULL,
		type				attendance_type NOT NULL,
		created_at			BIGINT NOT NULL,
		approved_at			BIGINT NOT NULL,
		approved_by_user_id	CHAR(36) NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (user_id, company_id, working_date)
	)`

	log.Println("CREATING attendance TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

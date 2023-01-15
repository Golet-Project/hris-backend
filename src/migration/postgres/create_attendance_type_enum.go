package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAttendanceTypeEnum(ctx context.Context) error {
	var sql = `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'attendance_type') THEN
			CREATE TYPE attendance_type AS ENUM(
				'schedule_change', 'work_time'
			);
		END IF;
	END$$;`

	log.Println("CREATING attendance_type ENUM")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

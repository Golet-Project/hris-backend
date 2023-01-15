package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAbsentTypeEnum(ctx context.Context) error {
	var sql = `
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'absent_type') THEN
			CREATE TYPE absent_type AS ENUM(
				'pto', 'truancy'
			);
		END IF;
	END$$;`

	log.Println("CREATING absent_type ENUM")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

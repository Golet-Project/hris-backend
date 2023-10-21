package tenant

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (m *Migration) CreateTableEmployee(ctx context.Context, tx pgx.Tx) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS public.employee (
		id							BIGSERIAL NOT NULL,
		uid							UUID NOT NULL DEFAULT uuid_generate_v4(),
		email						VARCHAR(50) NOT NULL,
		full_name				VARCHAR(255) NOT NULL,
		gender					VARCHAR(1) NOT NULL,
		birth_date			DATE NOT NULL,
		profile_picture	TEXT NOT NULL DEFAULT '',
		address					TEXT NOT NULL DEFAULT '',
		province_id			INTEGER NOT NULL DEFAULT 0,
		regency_id			INTEGER NOT NULL DEFAULT 0,
		district_id			INTEGER NOT NULL DEFAULT 0,
		join_date				DATE NOT NULL,
		end_date				DATE DEFAULT NULL,
		employee_status VARCHAR(20) NOT NULL,
		village_id			INTEGER NOT NULL DEFAULT 0,
		created_at			TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at			TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		UNIQUE (uid),
		UNIQUE (email)
	)`

	log.Println("CREATING employee TABLE")

	_, err := tx.Exec(ctx, sql)

	return err
}

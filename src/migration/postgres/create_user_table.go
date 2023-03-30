package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateUserTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS public.user (
		pk_id			BIGSERIAL NOT NULL,
		id				CHAR(36) NOT NULL,
		email			VARCHAR(50) NOT NULL,
		full_name		VARCHAR(255) NOT NULL,
		birth_date		DATE NOT NULL,
		profile_picture	TEXT NOT NULL,
		address			TEXT NOT NULL,
		province_id		INTEGER NOT NULL,
		regency_id		INTEGER NOT NULL,
		district_id		INTEGER NOT NULL,
		village_id		INTEGER NOT NULL,
		created_at		BIGINT NOT NULL,
		updated_at		BIGINT NOT NULL,

		PRIMARY KEY (pk_id),
		UNIQUE (id),
		UNIQUE (email)
	)`

	log.Println("CREATING user TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

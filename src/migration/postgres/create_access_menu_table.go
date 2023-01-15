package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateAccessMenuTable(ctx context.Context) error {
	var sql = `
	CREATE TABLE IF NOT EXISTS access_menu (
		id		SERIAL NOT NULL,
		name	VARCHAR(255) NOT NULL,

		PRIMARY KEY (id)
	)`

	log.Println("CREATING access_menu TABLE")

	_, err := m.Tx.Exec(ctx, sql)
	return err
}

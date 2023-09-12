package tenant

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (m *Migration) CreateExtensionUuid(ctx context.Context, tx pgx.Tx) error {
	var sql = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	_, err := tx.Exec(ctx, sql)
	if err != nil {
		log.Println("error when creating UUID EXTENSION")
		return err
	}

	return nil
}

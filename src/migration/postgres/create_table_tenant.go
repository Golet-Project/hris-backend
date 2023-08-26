package postgres

import (
	"context"
	"log"
)

func (m *Migrate) CreateTableTenant(ctx context.Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS tenant (
		id SERIAL NOT NULL,
		uid UUID NOT NULL DEFAULT uuid_generate_v4(),
		name VARCHAR(100) NOT NULL,
		domain VARCHAR(50) NOT NULL,

		PRIMARY KEY (id)
	);
	CREATE UNIQUE INDEX IF NOT EXISTS tenant_domain_key ON tenant(domain);
	CREATE UNIQUE INDEX IF NOT EXISTS tenant_uid_key ON tenant(uid);`

	_, err := m.Tx.Exec(ctx, sql)
	log.Println("CREATING TABLE tenant")

	return err
}

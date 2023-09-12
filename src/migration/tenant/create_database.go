package tenant

import (
	"context"
	"fmt"
	"log"
)

func (m *Migration) CreateDatabase(ctx context.Context, domain string) (dbName string, err error) {
	dbName = "tenant_" + domain
	var sql = fmt.Sprintf("CREATE DATABASE %s", dbName)

	log.Printf("CREATING tenant_%s DATABASE", domain)

	_, err = m.workerDBConn.Exec(ctx, sql)

	return
}

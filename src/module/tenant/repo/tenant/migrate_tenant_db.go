package tenant

import (
	"context"
	"fmt"
	"log"
	"strings"
)

func (r *Repository) MigrateTenantDB(ctx context.Context, domain string) error {
	log.Println("creating tenant db")
	dbName := fmt.Sprintf("tenant_%s", domain)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("CREATE DATABASE %s;", dbName))

	// TODO: Make base tables migration for tenant
	_, err := r.DB.Exec(ctx, builder.String())

	return err
}

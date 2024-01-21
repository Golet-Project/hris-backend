package queue

import (
	"context"
	"encoding/json"
	"hroost/infrastructure/store/postgres"
	"hroost/migration/tenant"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MigrateTenantDBIn struct {
	Domain string `json:"domain"`
}

func (q *Queue) MigrateTenantDB(ctx context.Context, in MigrateTenantDBIn) error {
	log.Println("migrating tenant DB")

	json, err := json.Marshal(in)
	if err != nil {
		return err
	}

	createdTask := asynq.NewTask(MigrateTenantDb, json)

	_, err = q.client.Enqueue(createdTask)
	if err != nil {
		return err
	}

	return nil
}

type MigrateTenantDbHandler struct {
	workerDBConn *pgxpool.Pool
	pgResolver   *postgres.Resolver
}

func NewMigrateTenantDbHandler(workerDBConn *pgxpool.Pool, pgResolver *postgres.Resolver) *MigrateTenantDbHandler {
	return &MigrateTenantDbHandler{
		workerDBConn: workerDBConn,
		pgResolver:   pgResolver,
	}
}

type TenantCreatePayload struct {
	Domain string
}

func (handler *MigrateTenantDbHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p TenantCreatePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	log.Printf("creating tenant database with domain: %s", p.Domain)

	migration, err := tenant.NewMigration(ctx, handler.workerDBConn, handler.pgResolver)
	if err != nil {
		return err
	}

	err = migration.Run(ctx, tenant.Tenant{
		Domain: p.Domain,
	})
	if err != nil {
		return err
	}

	return nil
}

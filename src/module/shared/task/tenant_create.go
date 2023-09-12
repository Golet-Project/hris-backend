package task

import (
	"context"
	"encoding/json"
	"hris/migration/tenant"
	"hris/module/shared/postgres"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

// === broadcaster ===
type TenantCreateTask struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

func NewTenantCreateTask(in TenantCreateTask) (*asynq.Task, error) {
	json, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TenantCreate, json), nil
}

// === handler ===
type TenantCreateHandler struct {
	workerDBConn *pgxpool.Pool
	pgResolver   *postgres.Resolver
}

func NewTenantCreateHandler(workerDBConn *pgxpool.Pool, pgResolver *postgres.Resolver) *TenantCreateHandler {
	return &TenantCreateHandler{
		workerDBConn: workerDBConn,
		pgResolver:   pgResolver,
	}
}

type TenantCreatePayload struct {
	Domain string
}

func (handler *TenantCreateHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
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

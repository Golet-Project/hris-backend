package main

import (
	"context"
	"fmt"
	"hroost/module/shared/postgres"
	tenantCentralTask "hroost/module/tenant/central/task"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	Mux    *asynq.ServeMux
	Server *asynq.Server
}

func NewWorker(workerDBConn *pgxpool.Pool, pgResolver *postgres.Resolver) (*Worker, error) {
	var redisUri = fmt.Sprintf("redis://%s@%s:%s/%s",
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_TASK_DB"),
	)

	clientOpt, err := asynq.ParseRedisURI(redisUri)
	if err != nil {
		return nil, err
	}

	var worker Worker

	worker.Server = asynq.NewServer(
		clientOpt,
		asynq.Config{
			Concurrency:  3,
			ErrorHandler: asynq.ErrorHandlerFunc(reportError),
		},
	)

	worker.Mux = asynq.NewServeMux()

	worker.Mux.Handle(tenantCentralTask.TenantCreate, tenantCentralTask.NewTenantCreateHandler(workerDBConn, pgResolver))

	return &worker, nil
}

func reportError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		err = fmt.Errorf("retry exhausted for task %s: %w", task.Type(), err)
	}

	// TODO: properly send to logger
	log.Println(err)
}

package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type Config struct {
	AsynqRedisMasterHost     string
	AsynqRedisMasterPort     string
	AsynqRedisMasterPassword string
	AsynqRedisMasterDb       int
}

type Worker struct {
	redisUri string

	mux    *asynq.ServeMux
	server *asynq.Server
}

func NewWorker(cfg *Config) (*Worker, error) {
	if cfg == nil {
		return nil, fmt.Errorf("[x] worker config required")
	}

	redisUri := fmt.Sprintf("redis://%s@%s:%s/%d",
		cfg.AsynqRedisMasterPassword,
		cfg.AsynqRedisMasterHost,
		cfg.AsynqRedisMasterPort,
		cfg.AsynqRedisMasterDb,
	)

	return &Worker{
		redisUri: redisUri,
	}, nil
}

// Run is blocking
func (w *Worker) Run(ctx context.Context) error {
	clientOpt, err := asynq.ParseRedisURI(w.redisUri)
	if err != nil {
		return err
	}

	w.server = asynq.NewServer(
		clientOpt,
		asynq.Config{
			Concurrency:  3,
			ErrorHandler: asynq.ErrorHandlerFunc(reportError),
		},
	)

	w.mux = asynq.NewServeMux()

	// TODO: register task

	log.Println("worker is running...")

	return w.server.Run(w.mux)
}

func (w *Worker) ShutDown(ctx context.Context) {
	w.server.Stop()
	w.server.Shutdown()
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

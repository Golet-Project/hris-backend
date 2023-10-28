package internal

import (
	"hris/module/tenant/internal/db"
	"hris/module/tenant/internal/queue"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Internal struct {
	db    *db.Db
	queue *queue.Queue
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	Queue      *asynq.Client
}

func New(d *Dependency) *Internal {
	if d.MasterConn == nil {
		log.Fatal("[x] Database connection required on tenant module")
	}
	if d.Queue == nil {
		log.Fatal("[x] Queue client required on tenant module")
	}

	dbImpl := db.New(&db.Dependency{
		MasterConn: d.MasterConn,
	})
	queueImpl := queue.New(&queue.Dependency{
		Client: d.Queue,
	})

	return &Internal{
		db:    dbImpl,
		queue: queueImpl,
	}
}

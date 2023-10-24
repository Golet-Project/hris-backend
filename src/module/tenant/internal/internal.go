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
	Pg    *pgxpool.Pool
	Queue *asynq.Client
}

func New(d *Dependency) *Internal {
	if d.Pg == nil {
		log.Fatal("[x] Database connection required on tenant module")
	}
	if d.Queue == nil {
		log.Fatal("[x] Queue client required on tenant module")
	}

	return &Internal{
		db: &db.Db{
			Pg: d.Pg,
		},
		queue: &queue.Queue{
			Client: d.Queue,
		},
	}
}

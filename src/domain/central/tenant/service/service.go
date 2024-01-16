package service

import (
	"fmt"
	"hroost/domain/central/tenant/db"
	"hroost/domain/central/tenant/queue"
)

type Config struct {
	Db    *db.Db
	Queue *queue.Queue
}

type Service struct {
	db    *db.Db
	queue *queue.Queue
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("db layer required")
	}
	if cfg.Queue == nil {
		return nil, fmt.Errorf("queue required")
	}

	return &Service{
		db:    cfg.Db,
		queue: cfg.Queue,
	}, nil
}

package service

import (
	"fmt"
	"hroost/central/domain/tenant/db"
	"hroost/central/domain/tenant/queue"
)

type Config struct {
	Db    db.IDbStore
	Queue queue.IQueue
}

type Service struct {
	db    db.IDbStore
	queue queue.IQueue
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

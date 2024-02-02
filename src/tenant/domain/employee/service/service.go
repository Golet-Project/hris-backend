package service

import (
	"fmt"
	"hroost/tenant/domain/employee/db"
)

type Config struct {
	Db db.IDbStore
}

type Service struct {
	db db.IDbStore
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("db layer required")
	}

	return &Service{
		db: cfg.Db,
	}, nil
}

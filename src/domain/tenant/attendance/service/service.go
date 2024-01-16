package service

import (
	"fmt"
	userService "hroost/domain/shared/user/service"
	"hroost/domain/tenant/attendance/db"
)

type Service struct {
	db *db.Db

	userService *userService.Service
}

type Config struct {
	Db *db.Db

	UserService *userService.Service
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("db required")
	}
	if cfg.UserService == nil {
		return nil, fmt.Errorf("userService required")
	}

	return &Service{
		db:          cfg.Db,
		userService: cfg.UserService,
	}, nil
}

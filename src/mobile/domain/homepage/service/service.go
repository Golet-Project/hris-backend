package service

import (
	"fmt"
	"hroost/mobile/domain/homepage/db"
	userService "hroost/shared/domain/user/service"
)

type Config struct {
	Db *db.Db

	// other service
	UserService *userService.Service
}

type Service struct {
	db *db.Db

	userService *userService.Service
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

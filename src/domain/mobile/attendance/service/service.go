package service

import (
	"fmt"

	"hroost/domain/mobile/attendance/db"
	userService "hroost/domain/shared/user/service"
)

type Service struct {
	db *db.Db

	// other service
	userService *userService.Service
}

type Config struct {
	Db *db.Db

	// other service
	UserService *userService.Service
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("Db layer required")
	}

	if cfg.UserService == nil {
		return nil, fmt.Errorf("userService required")
	}

	return &Service{
		db:          cfg.Db,
		userService: cfg.UserService,
	}, nil
}

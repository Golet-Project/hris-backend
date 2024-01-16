package service

import (
	"fmt"
	"hroost/domain/mobile/auth/db"
	"hroost/domain/mobile/auth/memory"

	userService "hroost/domain/shared/user/service"
)

type Service struct {
	db     *db.Db
	memory *memory.Memory

	// other service
	userService *userService.Service
}

type Config struct {
	Db     *db.Db
	Memory *memory.Memory

	// other service
	UserService *userService.Service
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Db == nil {
		return nil, fmt.Errorf("db required")
	}
	if cfg.Memory == nil {
		return nil, fmt.Errorf("memory required")
	}
	if cfg.UserService == nil {
		return nil, fmt.Errorf("userService required")
	}

	return &Service{
		db:          cfg.Db,
		memory:      cfg.Memory,
		userService: cfg.UserService,
	}, nil
}

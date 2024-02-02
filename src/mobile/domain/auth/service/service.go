package service

import (
	"fmt"
	"hroost/mobile/domain/auth/db"
	"hroost/mobile/domain/auth/memory"

	userService "hroost/shared/domain/user/service"
)

type Service struct {
	db     db.IDbStore
	memory memory.IMemory

	// other service
	userService *userService.Service
}

type Config struct {
	Db     db.IDbStore
	Memory memory.IMemory

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

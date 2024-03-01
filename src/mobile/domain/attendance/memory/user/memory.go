package user

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis *redis.Client
}

type Memory struct {
	redis *redis.Client
}

func New(cfg *Config) (*Memory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Redis == nil {
		return nil, fmt.Errorf("redis required")
	}

	return &Memory{
		redis: cfg.Redis,
	}, nil
}

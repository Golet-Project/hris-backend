package memory

import (
	"fmt"

	redisClient "github.com/redis/go-redis/v9"
)

type Config struct {
	Client *redisClient.Client
}

type Memory struct {
	client *redisClient.Client
}

func New(cfg *Config) (*Memory, error) {
	if cfg.Client == nil {
		return nil, fmt.Errorf("config required for memory")
	}

	return &Memory{
		client: cfg.Client,
	}, nil
}

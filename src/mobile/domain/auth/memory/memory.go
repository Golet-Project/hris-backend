package memory

import (
	"context"
	"fmt"

	redisClient "github.com/redis/go-redis/v9"
)

type IMemory interface {
	DeletePasswordRecoveryToken(ctx context.Context, userId string) (err error)

	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err error)

	SetPasswordRecoveryToken(ctx context.Context, userId, token string) (err error)
}

type Config struct {
	Client *redisClient.Client
}

type Memory struct {
	client *redisClient.Client
}

func New(cfg *Config) (*Memory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config required")
	}
	if cfg.Client == nil {
		return nil, fmt.Errorf("client required")
	}

	return &Memory{
		client: cfg.Client,
	}, nil
}

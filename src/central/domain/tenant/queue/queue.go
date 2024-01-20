package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
)

type Queue struct {
	client *asynq.Client
}

type Config struct {
	Client *asynq.Client
}

func New(cfg *Config) (*Queue, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config requried")
	}
	if cfg.Client == nil {
		return nil, fmt.Errorf("client required")
	}

	return &Queue{
		client: cfg.Client,
	}, nil
}

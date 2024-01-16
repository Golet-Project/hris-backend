package redis

import (
	"fmt"

	"github.com/hibiken/asynq"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Db       int
}

type redis struct {
	redisAddr string
	db        int
	password  string
}

func NewRedis(cfg *RedisConfig) (*redis, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis queue config is required")
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &redis{
		redisAddr: addr,
		db:        cfg.Db,
		password:  cfg.Password,
	}, nil
}

func (s *redis) Create() *asynq.Client {
	queueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     s.redisAddr,
		DB:       s.db,
		Password: s.password,
	})

	return queueClient
}

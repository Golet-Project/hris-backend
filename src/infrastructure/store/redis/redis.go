package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Db       int
}

type redisDb struct {
	redisAddr string
	db        int
	password  string
}

func NewRedis(cfg *RedisConfig) (*redisDb, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis config required")
	}

	redisAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &redisDb{
		redisAddr: redisAddr,
		db:        cfg.Db,
	}, nil
}

func (r *redisDb) Connect(ctx context.Context) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.redisAddr,
		Password: r.password,
		DB:       r.db,

		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("[v] Redis connected...")
			return nil
		},
	})
}

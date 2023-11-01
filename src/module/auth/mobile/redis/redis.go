package redis

import (
	"log"

	redisClient "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redisClient.Client
}

type Dependency struct {
	Client *redisClient.Client
}

func New(d *Dependency) *Redis {
	if d.Client == nil {
		log.Fatal("[x] Redis connection required on auth/mobile/redis module")
	}

	return &Redis{
		client: d.Client,
	}
}

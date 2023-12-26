package redis

import redisClient "github.com/redis/go-redis/v9"

type Redis struct {
	client *redisClient.Client
}
type Dependency struct {
	Client *redisClient.Client
}

func New(d *Dependency) *Redis {
	if d.Client == nil {
		panic("[x] Redis connection required on auth/central/redis module")
	}

	return &Redis{
		client: d.Client,
	}
}

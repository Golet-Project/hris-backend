package redis

import redisClient "github.com/redis/go-redis/v9"

type Redis struct {
	Client *redisClient.Client
}

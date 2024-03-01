package user

import (
	"context"
	"errors"
	"hroost/shared/primitive"

	"github.com/redis/go-redis/v9"
)

const domainRedisKey = "domain:"

func (m *Memory) GetDomain(ctx context.Context, userId string) (domain string, repoError *primitive.RepoError) {
	domain, err := m.redis.Get(ctx, domainRedisKey+userId).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return domain, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
				Err:   err,
			}
		}
	}

	return domain, nil
}

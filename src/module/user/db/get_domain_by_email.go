package db

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const domainRedisKey = "domain:"

func (d *Db) GetDomainByEmail(ctx context.Context, email string) (domain string, err error) {
	// get from redis first
	domain, err = d.redis.Get(ctx, domainRedisKey+email).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return domain, err
		}
	} else {
		return domain, nil
	}

	sql := `SELECT domain FROM users WHERE email = $1 AND deleted_at IS NULL`

	err = d.masterConn.QueryRow(ctx, sql, email).Scan(&domain)
	if err != nil {
		return domain, err
	}

	// set to redis
	err = d.redis.Set(ctx, domainRedisKey+email, domain, time.Hour*24).Err()

	return
}

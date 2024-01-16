package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"time"

	"github.com/redis/go-redis/v9"
)

func (d *Db) GetDomainByUid(ctx context.Context, uid string) (domain string, err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	// get from redis
	domain, err = d.redis.Get(ctx, domainRedisKey+uid).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return domain, err
		}
	} else {
		return domain, nil
	}

	var sql = `
	SELECT
		domain
	FROM
		users
	WHERE
		uid = $1`

	err = masterConn.QueryRow(ctx, sql, uid).Scan(&domain)
	if err != nil {
		return domain, err
	}

	// set to redis
	err = d.redis.Set(ctx, domainRedisKey+uid, domain, time.Hour*24).Err()

	return
}

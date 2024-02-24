package db

import (
	"context"
	"errors"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const domainRedisKey = "domain:"

func (d *Db) GetDomainByUid(ctx context.Context, uid string) (domain string, repoError *primitive.RepoError) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return domain, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	// TODO: separate redis module
	// get from redis
	domain, err = d.redis.Get(ctx, domainRedisKey+uid).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return domain, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
				Err:   err,
			}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
				Err:   err,
			}
		}

		return domain, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	// set to redis
	err = d.redis.Set(ctx, domainRedisKey+uid, domain, time.Hour*24).Err()
	if err != nil {
		return domain, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}
	}

	return
}

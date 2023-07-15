package auth

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) RedisGetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	token, err = r.Redis.Get(ctx, key).Result()

	return
}

func (r *Repository) RedisSetPasswordRecoveryToken(ctx context.Context, userId, token string) error {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	_, err := r.Redis.SetEx(ctx, key, token, time.Minute * 3).Result()

	return err
}

func (r *Repository) RedisDeletePasswordRecoveryToken(ctx context.Context, userId string) error {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	_, err := r.Redis.Del(ctx, key).Result()

	return err
}
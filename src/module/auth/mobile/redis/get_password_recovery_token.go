package redis

import (
	"context"
	"fmt"
)

func (r *Redis) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	token, err = r.client.Get(ctx, key).Result()

	return
}

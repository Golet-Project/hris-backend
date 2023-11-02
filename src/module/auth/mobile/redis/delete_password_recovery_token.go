package redis

import (
	"context"
	"fmt"
)

func (r *Redis) DeletePasswordRecoveryToken(ctx context.Context, userId string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = r.client.Del(ctx, key).Err()

	return
}

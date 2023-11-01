package redis

import (
	"context"
	"fmt"
	"time"
)

func (r *Redis) SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = r.client.Set(ctx, key, token, time.Minute*3).Err()

	return
}

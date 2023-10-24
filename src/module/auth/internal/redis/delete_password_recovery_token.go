package redis

import (
	"context"
	"fmt"
)

// DeletePasswordRecoveryToken delete the user password recovery token
func (r *Redis) DeletePasswordRecoveryToken(ctx context.Context, userId string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = r.Client.Del(ctx, key).Err()

	return
}
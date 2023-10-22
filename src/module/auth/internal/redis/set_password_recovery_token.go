package redis

import (
	"context"
	"fmt"
	"time"
)

// SetPasswordRecoveryToken set the user password recovery token
// within the specified expiration time
func (r *Redis) SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = r.Client.Set(ctx, key, token, time.Minute * 3).Err()

	return
}

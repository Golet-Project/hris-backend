package memory

import (
	"context"
	"fmt"
)

// DeletePasswordRecoveryToken delete the user password recovery token
func (m *Memory) DeletePasswordRecoveryToken(ctx context.Context, userId string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = m.client.Del(ctx, key).Err()

	return
}

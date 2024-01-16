package memory

import (
	"context"
	"fmt"
)

func (m *Memory) DeletePasswordRecoveryToken(ctx context.Context, userId string) (err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err = m.client.Del(ctx, key).Err()

	return
}

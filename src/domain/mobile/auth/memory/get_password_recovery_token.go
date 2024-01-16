package memory

import (
	"context"
	"fmt"
)

func (m *Memory) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err error) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	token, err = m.client.Get(ctx, key).Result()

	return
}

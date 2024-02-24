package memory

import (
	"context"
	"fmt"
	"hroost/shared/primitive"
	"time"
)

func (m *Memory) SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (repoError *primitive.RepoError) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err := m.client.Set(ctx, key, token, time.Minute*3).Err()
	if err != nil {
		return &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	return
}

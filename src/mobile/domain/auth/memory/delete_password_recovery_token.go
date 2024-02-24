package memory

import (
	"context"
	"fmt"
	"hroost/shared/primitive"
)

func (m *Memory) DeletePasswordRecoveryToken(ctx context.Context, userId string) (repoError *primitive.RepoError) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err := m.client.Del(ctx, key).Err()
	if err != nil {
		return &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	return
}

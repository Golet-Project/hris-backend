package memory

import (
	"context"
	"fmt"
	"hroost/shared/primitive"
	"time"
)

// SetPasswordRecoveryToken set the user password recovery token
// within the specified expiration time
func (m *Memory) SetPasswordRecoveryToken(ctx context.Context, userId string, token string) *primitive.RepoError {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	err := m.client.Set(ctx, key, token, time.Minute*3).Err()
	if err != nil {
		return &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
		}
	}

	return nil
}

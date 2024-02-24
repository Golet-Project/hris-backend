package memory

import (
	"context"
	"errors"
	"fmt"
	"hroost/shared/primitive"

	"github.com/redis/go-redis/v9"
)

func (m *Memory) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, repoError *primitive.RepoError) {
	key := fmt.Sprintf("password_recovery_token_%s", userId)

	token, err := m.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
			}
		} else {
			return "", &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
			}
		}
	}

	return token, nil
}

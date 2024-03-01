package user_test

import (
	"hroost/mobile/domain/attendance/memory/user"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		_, err := user.New(nil)
		if err == nil {
			t.Error("expect error got nil")
		}
	})

	redisClient := redis.NewClient(&redis.Options{})

	t.Run("valid param", func(t *testing.T) {
		uMem, err := user.New(&user.Config{
			Redis: redisClient,
		})

		t.Run("nil error", func(t *testing.T) {
			if err != nil {
				t.Error("expect nil got error")
			}
		})

		t.Run("valid instance", func(t *testing.T) {
			if uMem == nil {
				t.Errorf("expect %T got nil", user.Memory{})
			}
		})
	})

	t.Run("missing Redis", func(t *testing.T) {
		_, err := user.New(&user.Config{})
		if err == nil {
			t.Error("expect error but got nil")
		}

		var errMsg = "redis required"
		if err.Error() != errMsg {
			t.Errorf("expect error message: '%s' but got '%s'", errMsg, err.Error())
		}
	})
}

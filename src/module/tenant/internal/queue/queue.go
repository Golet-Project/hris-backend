package queue

import "github.com/hibiken/asynq"

type Queue struct {
	Client *asynq.Client
}

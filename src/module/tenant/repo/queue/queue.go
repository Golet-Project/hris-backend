package queue

import "github.com/hibiken/asynq"

type QueueRepo struct {
	Client *asynq.Client
}

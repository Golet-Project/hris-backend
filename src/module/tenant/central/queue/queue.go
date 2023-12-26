package queue

import "github.com/hibiken/asynq"

type Queue struct {
	client *asynq.Client
}

type Dependency struct {
	Client *asynq.Client
}

func New(d *Dependency) *Queue {
	if d.Client == nil {
		panic("[x] Queue client required on tenant/central/queue module")
	}

	return &Queue{
		client: d.Client,
	}
}

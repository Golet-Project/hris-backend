package queue

import (
	"context"
	"hris/module/shared/task"
	"log"
)

type MigrateTenantDBIn struct {
	Domain string
}

func (q *QueueRepo) MigrateTenantDB(ctx context.Context, in MigrateTenantDBIn) error {
	log.Println("creating tenant DB")

	createdTask, err := task.NewTenantCreateTask(task.TenantCreateTask{
		Domain: in.Domain,
	})
	if err != nil {
		return err
	}

	_, err = q.Client.Enqueue(createdTask)
	if err != nil {
		return err
	}

	return nil
}

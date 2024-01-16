package queue

import (
	"context"
	"hroost/domain/central/tenant/task"
	"log"
)

type MigrateTenantDBIn struct {
	Domain string
}

func (q *Queue) MigrateTenantDB(ctx context.Context, in MigrateTenantDBIn) error {
	log.Println("creating tenant DB")

	createdTask, err := task.NewTenantCreateTask(task.TenantCreateTask{
		Domain: in.Domain,
	})
	if err != nil {
		return err
	}

	_, err = q.client.Enqueue(createdTask)
	if err != nil {
		return err
	}

	return nil
}

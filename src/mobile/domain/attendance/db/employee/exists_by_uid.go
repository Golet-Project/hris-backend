package employee

import (
	"context"
	"hroost/shared/primitive"
)

func (d *Db) ExistsById(ctx context.Context, domain string, uid string) (exist bool, repoError *primitive.RepoError) {
	panic("implement me")
}

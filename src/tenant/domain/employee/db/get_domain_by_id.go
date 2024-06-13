package db

import (
	"context"
	"errors"
	"fmt"
	"hroost/infrastructure/store/postgres"
	"hroost/shared/primitive"

	"github.com/jackc/pgx/v5"
)

func (d *Db) GetDomainById(ctx context.Context, id string) (domain string, repoErr *primitive.RepoError) {
	if id == "" {
		repoErr = &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("employee domain not found"),
		}
	}

	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		repoErr = &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   err,
		}

		return
	}

	sql := `SELECT domain FROM users WHERE uid = $1`

	err = masterConn.QueryRow(ctx, sql, id).Scan(&domain)
	if err != nil {
		var issue primitive.RepoErrorCode = primitive.RepoErrorCodeServerError

		if errors.Is(err, pgx.ErrNoRows) {
			issue = primitive.RepoErrorCodeDataNotFound
		}

		repoErr = &primitive.RepoError{
			Issue: issue,
			Err:   err,
		}

		return
	}

	return
}

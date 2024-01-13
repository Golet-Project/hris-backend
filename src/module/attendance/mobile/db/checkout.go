package db

import (
	"context"
	"hroost/module/shared/postgres"
)

func (d *Db) Checkout(ctx context.Context, domain, uid string) (rowsAffected int64, err error) {
	var sql = `
	UPDATE
		attendance
	SET
		checkout_time = now()
	WHERE
		employee_uid = $1
		AND
		checkout_time IS NULL`

	conn, err := d.pgResolver.Resolve(postgres.Domain(domain))
	if err != nil {
		return
	}

	commandTag, err := conn.Exec(ctx, sql, uid)

	rowsAffected = commandTag.RowsAffected()

	return
}

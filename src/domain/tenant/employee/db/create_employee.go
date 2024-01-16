package db

import (
	"context"
	"fmt"
	"hroost/infrastructure/store/postgres"
	"hroost/module/shared/primitive"
)

type CreateEmployeeIn struct {
	Domain         string
	Email          string
	Password       string
	FirstName      string
	LastName       string
	Gender         primitive.Gender
	BirthDate      string
	Address        string
	ProvinceId     string
	RegencyId      string
	DistrictId     string
	VillageId      string
	JoinDate       string
	EmployeeStatus primitive.EmployeeStatus
}

func (d *Db) CreateEmployee(ctx context.Context, data CreateEmployeeIn) (err error) {
	masterConn, err := d.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		return
	}

	// insert into master
	masterTx, err := masterConn.Begin(ctx)
	if err != nil {
		return
	}
	var insertToMasterSql = `
	INSERT INTO users (domain, email, password, first_name, last_name, birth_date)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING uid`

	var insertedUid string
	err = masterTx.QueryRow(ctx, insertToMasterSql,
		data.Domain, data.Email, data.FirstName, data.LastName, data.BirthDate,
	).Scan(&insertedUid)
	if err != nil {
		if err2 := masterTx.Rollback(ctx); err2 != nil {
			return err2
		}

		return
	}

	// insert into tenant
	var insertToTenantSql = `
	INSERT INTO employee (
		uid, email, full_name, gender, birth_date, address, province_id, regency_id,
		district_id, village_id, join_date, employee_status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	tenantConn, err := d.pgResolver.Resolve(postgres.Domain(data.Domain))
	if err != nil {
		if err2 := masterTx.Rollback(ctx); err2 != nil {
			return err2
		}

		return
	}

	_, err = tenantConn.Exec(ctx, insertToTenantSql,
		insertedUid, data.Email, fmt.Sprintf("%s %s", data.FirstName, data.LastName), data.Gender, data.BirthDate, data.Address,
		data.ProvinceId, data.RegencyId, data.DistrictId, data.VillageId, data.JoinDate, data.EmployeeStatus,
	)
	if err != nil {
		if err2 := masterTx.Rollback(ctx); err2 != nil {
			return err2
		}

		return
	}

	err = masterTx.Commit(ctx)
	if err != nil {
		return
	}

	return
}

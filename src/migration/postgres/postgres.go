package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Migrate struct {
	Tx pgx.Tx
}

func (m *Migrate) RunMigration(ctx context.Context) error {
	// if err := m.CreateAttendanceTypeEnum(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateAbsentTypeEnum(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateAbsentQuotaTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateAbsentTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateAccessMenuTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateAttendanceTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateCompanyTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateRoleAccessMenuTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateRoleTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateUserAccessMenuTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateUserCompanyTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateUserCompanyWorkTimeTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateUserRoleAccessMenuTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateUserRoleTable(ctx); err != nil {
	// 	return err
	// }

	if err := m.CreateEmployeeTable(ctx); err != nil {
		return err
	}

	// if err := m.CreateWorkTimeTable(ctx); err != nil {
	// 	return err
	// }

	return nil
}

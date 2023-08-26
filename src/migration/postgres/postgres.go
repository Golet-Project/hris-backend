package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Migrate struct {
	Tx pgx.Tx
}

func (m *Migrate) RunMigration(ctx context.Context) error {
	// NOTE: for development, we'll comment on this first

	// if err := m.CreateEmployeeTable(ctx); err != nil {
	// 	return err
	// }

	// if err := m.SeedWilayah(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateTableAndSeedProvince(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateTableAndSeedRegency(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateTableAndSeedDistrict(ctx); err != nil {
	// 	return err
	// }

	// if err := m.CreateTableAndSeedVillage(ctx); err != nil {
	// 	return err
	// }

	if err := m.CreateTableTenant(ctx); err != nil {
		return err
	}

	return nil
}

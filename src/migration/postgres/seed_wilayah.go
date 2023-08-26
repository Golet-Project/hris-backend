package postgres

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (m *Migrate) SeedWilayah(ctx context.Context) error {
	path := filepath.Join("./postgres", "wilayah.sql")

	c, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		return ioErr
	}

	log.Println("SEEDING wilayah")

	sql := string(c)
	_, err := m.Tx.Exec(ctx, sql)
	return err
}

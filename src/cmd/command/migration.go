package command

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
)

func migrationCommand() *cli.Command {
	return &cli.Command{
		Name:     "migrate",
		Usage:    "command about database migration",
		Category: "Migration",
		Subcommands: []*cli.Command{
			{
				Name:  "up",
				Usage: "run database migration",
				Action: func(cCtx *cli.Context) error {
					databaseURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s",
						os.Getenv("PG_MASTER_USER"),
						os.Getenv("PG_MASTER_PASSWORD"),
						os.Getenv("PG_MASTER_HOST"),
						os.Getenv("PG_MASTER_PORT"),
						os.Getenv("PG_MASTER_DATABASE"),
					)

					log.Println("running migration up...")
					m, err := migrate.New(
						"file://migration/master/postgres",
						databaseURL,
					)
					if err != nil {
						return err
					}
					defer m.Close()

					err = m.Up()
					if err != nil {
						return err
					}

					log.Println("done!")

					return nil
				},
			},
			{
				Name:  "down",
				Usage: "rollback database migration",
				Action: func(cCtx *cli.Context) error {
					databaseURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s",
						os.Getenv("PG_MASTER_USER"),
						os.Getenv("PG_MASTER_PASSWORD"),
						os.Getenv("PG_MASTER_HOST"),
						os.Getenv("PG_MASTER_PORT"),
						os.Getenv("PG_MASTER_DATABASE"),
					)

					log.Println("running migration down...")
					m, err := migrate.New(
						"file://migration/master/postgres",
						databaseURL,
					)
					if err != nil {
						return err
					}
					defer m.Close()

					err = m.Down()
					if err != nil {
						return err
					}

					log.Println("done!")

					return nil
				},
			},
		},
	}
}

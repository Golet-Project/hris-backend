package main

import (
	"context"
	"fmt"
	"hroost/migration/postgres"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	var pgConnString string
	if val, ok := os.LookupEnv("PG_URL"); !ok {
		log.Fatal("must provide PG_URL environment variable")
	} else {
		pgConnString = val
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, pgConnString)
	if err != nil {
		fmt.Println("[x] Failed to connect postgreSQL")
		log.Fatal(err)
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	postgres := postgres.Migrate{
		Tx: tx,
	}

	err = postgres.RunMigration(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			log.Fatal(e)
		}
		log.Fatal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"context"
	"fmt"
	"hris/cmd"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// parse config
	cfg := parseConfig()

	// create the http server
	app := cmd.NewHttpServer(fiber.Config{
		AppName: cfg.appName,
	})

	// initialize database
	pgPool, err := pgxpool.New(context.Background(), cfg.pgUrl)
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
	}
	defer pgPool.Close()

	if err := pgPool.Ping(context.Background()); err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
	} else {
		fmt.Println("[v] PostgreSQL connected...")
	}

	// start the app
	var errChan = make(chan error, 1)
	go func() {
		err := app.Listen(cfg.httpPort)

		if err != nil {
			errChan <- err
		}
	}()

	log.Fatal(<-errChan)
}

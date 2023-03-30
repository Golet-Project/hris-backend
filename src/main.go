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

	// initialize database
	connConfig, err := pgxpool.ParseConfig(cfg.pgUrl)
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
	}
	connConfig.MinConns = 3
	connConfig.MaxConns = 5

	pgPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
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

	// create the http server
	app := cmd.NewApp(cmd.AppConfig{
		DB: pgPool,

		FiberCfg: fiber.Config{
			AppName: cfg.appName,
		},
	})

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

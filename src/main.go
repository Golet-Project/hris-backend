package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

	//=== INITIALIZE REDIS ===
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: redisDB,

		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("[v] Redis connected...")
			return nil
		},
	})

	// create the http server
	app := NewApp(AppConfig{
		DB: pgPool,
		Redis: rdb,

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

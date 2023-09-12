package main

import (
	"context"
	"errors"
	"fmt"
	"hris/module/shared/postgres"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	cfg := parseConfig()

	// init database
	pgPool := initDatabase(cfg)
	defer pgPool.Close()

	// register to resolver
	pgResolver := postgres.NewResolver(pgPool)
	tenants, err := getAllTenant(pgPool)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Fatal("[x] Error during get all tenant")
		}
	}
	for _, tenant := range tenants {
		// store to resolver
		tenantPool, err := makeTenantConnection(tenant.Domain)
		if err != nil {
			log.Fatal("[x] error during creating tenant connection", err)
		}

		pgResolver.Register(postgres.Database{
			DomainName: postgres.Domain(tenant.Domain),
			Pool:       tenantPool,
		})
	}

	// init worker db conn
	workerDBConn := initWorkerDatabaseConnn(cfg)
	defer workerDBConn.Close()

	// init queue client
	queueClient := initQueueClient(cfg)
	defer func() {
		if err := queueClient.Close(); err != nil {
			log.Println("[x] error during closing queue connection")
		}
	}()

	//=== INITIALIZE REDIS ===
	rdb := initRedis(cfg)

	// create the http server
	app := NewApp(AppConfig{
		DB:               pgPool,
		Redis:            rdb,
		QueueClient:      queueClient,
		PostgresResolver: pgResolver,

		FiberCfg: fiber.Config{
			AppName: cfg.appName,
		},
	})

	// create worker server
	worker, err := NewWorker(workerDBConn, pgResolver)
	if err != nil {
		log.Fatal(err)
		return
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// start the app
	var errChan = make(chan error, 1)
	go func() {
		err := app.Listen(cfg.httpPort)

		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		log.Println("worker is running...")
		err := worker.Server.Run(worker.Mux)

		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-sigChan:
		var shutdownWg sync.WaitGroup

		go func() {
			shutdownWg.Add(1)
			defer shutdownWg.Done()
			// loop trhough db connection
			log.Println("closing database resolver connection")
			masterConn := pgResolver.MasterConn
			masterConn.Close()

			allTenantConn := pgResolver.GetAllTenantConn()
			if allTenantConn != nil {
				allTenantConn.Range(func(key, val interface{}) bool {
					log.Printf("closing %s database connection\n", key)
					pool := val.(*pgxpool.Pool)
					pool.Close()
					return true
				})
			}
		}()

		go func() {
			shutdownWg.Add(1)
			defer shutdownWg.Done()

			log.Println("shutdown app...")
			if err := app.Shutdown(); err != nil {
				log.Fatal("failed to shutdown", err)
			} else {
				log.Println("application is shutdown properly")
			}
		}()
		shutdownWg.Wait()
		log.Println("done!")
	}
}

func makeTenantConnection(domain string) (conn *pgxpool.Pool, err error) {
	dbName := "tenant_" + domain
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		dbName,
	)
	connConfig, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		log.Printf("[x] Failed to make tenant: %s database connection\n", domain)
		return nil, err
	}
	connConfig.MinConns = 3
	connConfig.MaxConns = 5

	pgPool, err := pgxpool.NewWithConfig(context.TODO(), connConfig)
	if err != nil {
		return nil, err
	}

	if err := pgPool.Ping(context.TODO()); err != nil {
		log.Printf("[x] Failed to make tenant: %s database connection\n", domain)
		return nil, err
	}

	log.Printf("[v] tenant: %s database connected", domain)

	return pgPool, nil
}

type GetAllTenantOut struct {
	ID     string
	Name   string
	Domain string
}

func getAllTenant(masterConn *pgxpool.Pool) (out []GetAllTenantOut, err error) {
	var sql = `SELECT uid, name, domain FROM tenant`

	rows, err := masterConn.Query(context.Background(), sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tenant GetAllTenantOut

		err = rows.Scan(&tenant.ID, &tenant.Name, &tenant.Domain)
		if err != nil {
			return
		}

		out = append(out, tenant)
	}

	return
}

func initDatabase(cfg config) *pgxpool.Pool {
	// initialize database
	connConfig, err := postgres.MasterConnConfig()
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
		return nil
	}

	if err := pgPool.Ping(context.Background()); err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		log.Fatal(err)
		return nil
	} else {
		fmt.Println("[v] PostgreSQL connected...")
	}

	return pgPool
}

func initWorkerDatabaseConnn(cfg config) *pgxpool.Pool {
	// init worker conn
	workerDBConnConfig, err := postgres.WorkerConnConfig()
	if err != nil {
		log.Println("[x] Failed to make DB connection for worker process")
		log.Fatal(err)
	}

	workerDBConn, err := pgxpool.NewWithConfig(context.Background(), workerDBConnConfig)

	if err := workerDBConn.Ping(context.Background()); err != nil {
		fmt.Println("[x] Failed to connect worker PostgreSQL")
		log.Fatal(err)
		return nil
	} else {
		fmt.Println("[v] PostgreSQL for worker connected...")
	}

	return workerDBConn
}

func initQueueClient(cfg config) *asynq.Client {
	db, err := strconv.Atoi(cfg.redisTaskDB)
	if err != nil {
		log.Fatal("invalid REDIS_TASK_DB")
		return nil
	}

	var host = cfg.redisHost
	var port = cfg.redisPort

	var addr = fmt.Sprintf("%s:%s", host, port)

	queueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: addr,
		DB:   db,
	})

	return queueClient
}

func initRedis(cfg config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", cfg.redisHost, cfg.redisPort)
	redisDB, _ := strconv.Atoi(cfg.redisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.redisPassword,
		DB:       redisDB,

		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("[v] Redis connected...")
			return nil
		},
	})

	return rdb
}

package server

import (
	"context"
	"errors"
	"fmt"
	redisQueue "hroost/infrastructure/queue/redis"
	"hroost/infrastructure/store/postgres"
	redisDb "hroost/infrastructure/store/redis"
	"hroost/infrastructure/worker"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	cfg config

	dbConnFuncs []func(context.Context) error
	pgResolver  *postgres.Resolver

	queueClientFunc func() error
	queueClient     *asynq.Client

	redisClientFunc func(context.Context) error
	redisClient     *redis.Client

	worker *worker.Worker

	app *fiber.App
}

func NewServer() *Server {
	cfg := parseConfig()

	return &Server{
		cfg:        cfg,
		pgResolver: postgres.NewResolver(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// db
	err := s.withMasterDbConn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = s.withTenantDbConn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = s.withWorkerDbConn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// redis
	err = s.withRedis(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// queue
	err = s.withQueueClient()
	if err != nil {
		log.Fatal(err)
	}

	// worker
	err = s.withWorker(ctx)
	if err != nil {
		log.Fatal(err)
	}

	s.fiber()

	// run app
	g.Go(func() error {
		log.Println("app listening on port", s.cfg.httpPort)
		return s.app.Listen(s.cfg.httpPort)
	})

	g.Go(func() error {
		log.Println("runing worker")
		return s.worker.Run(ctx)
	})

	return g.Wait()
}

func (s *Server) withMasterDbConn(ctx context.Context) error {
	db, err := postgres.NewMasterDb(&postgres.MasterDbConfig{
		User:     s.cfg.pgMasterUser,
		Password: s.cfg.pgMasterPassword,
		Host:     s.cfg.pgMasterHost,
		Port:     s.cfg.pgMasterPort,
		Db:       s.cfg.pgMasterDatabase,
	})
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		return err
	}

	pgPool, err := db.Connect(ctx)
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		return err
	}

	err = s.pgResolver.Register(postgres.Database{
		DomainName: postgres.MasterDomain,
		Pool:       pgPool,
	})
	if err != nil {
		fmt.Println("[x] Failed to connect PostgreSQL")
		return err
	}

	fmt.Println("[v] PostgreSQL connected...")
	return nil
}

func (s *Server) withTenantDbConn(ctx context.Context) error {
	// get master conn
	var masterConn *pgxpool.Pool
	var err error

	masterConn, err = s.pgResolver.Resolve(postgres.MasterDomain)
	if err != nil {
		if !errors.Is(err, postgres.ErrNoConnection) {
			return err
		}

		// make temp conn if not exists
		db, err := postgres.NewMasterDb(&postgres.MasterDbConfig{
			User:     s.cfg.pgMasterUser,
			Password: s.cfg.pgMasterPassword,
			Host:     s.cfg.pgMasterHost,
			Port:     s.cfg.pgMasterPort,
			Db:       s.cfg.pgMasterDatabase,
		})
		if err != nil {
			fmt.Println("[x] Failed to connect PostgreSQL")
			return err
		}

		pgPool, err := db.Connect(ctx)
		if err != nil {
			fmt.Println("[x] Failed to connect PostgreSQL")
			return err
		}
		defer pgPool.Close()

		masterConn = pgPool

		fmt.Println("[v] PostgreSQL connected...")
	}

	// get all tenant from master db
	var sql = `SELECT uid, name, domain FROM tenant`
	rows, err := masterConn.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	type tenantStruct struct {
		uid    string
		name   string
		domain string
	}
	var tenants []tenantStruct

	for rows.Next() {
		var tenant tenantStruct

		err = rows.Scan(&tenant.uid, &tenant.name, &tenant.domain)
		if err != nil {
			return err
		}

		tenants = append(tenants, tenant)
	}

	// make a connection for each tenant
	g, ctx := errgroup.WithContext(ctx)
	for _, each := range tenants {
		go func(tenant tenantStruct) {
			g.Go(func() error {
				db, err := postgres.NewTenantDb(&postgres.TenantDbConfig{
					Domain: tenant.domain,
					// NOTE: karena masih 1 instance db, jadi pakai user master
					User:     s.cfg.pgMasterUser,
					Password: s.cfg.pgMasterPassword,
					Host:     s.cfg.pgMasterHost,
					Port:     s.cfg.pgMasterPort,
				})
				if err != nil {
					return err
				}

				pgPool, err := db.Connect(ctx)
				if err != nil {
					return err
				}

				err = s.pgResolver.Register(postgres.Database{
					DomainName: postgres.Domain(tenant.domain),
					Pool:       pgPool,
				})
				if err != nil {
					return err
				}

				log.Printf("[v] tenant: %s database connected", tenant.domain)

				return nil
			})
		}(each)
	}

	return g.Wait()
}

func (s *Server) withWorkerDbConn(ctx context.Context) error {
	db, err := postgres.NewWorkerDb(&postgres.WorkerDbConfig{
		// NOTE: sementara pakai master conn
		User:     s.cfg.pgMasterUser,
		Password: s.cfg.pgMasterPassword,
		Host:     s.cfg.pgMasterHost,
		Port:     s.cfg.pgMasterPort,
		Db:       s.cfg.pgMasterDatabase,
	})
	if err != nil {
		return err
	}

	pgPool, err := db.Connect(ctx)
	if err != nil {
		return err
	}

	err = s.pgResolver.Register(postgres.Database{
		DomainName: postgres.Domain("worker"),
		Pool:       pgPool,
	})
	if err != nil {
		return err
	}

	fmt.Println("[v] PostgreSQL for worker connected...")

	return nil
}

func (s *Server) withRedis(ctx context.Context) error {
	db, err := redisDb.NewRedis(&redisDb.RedisConfig{
		Host:     s.cfg.redisMasterHost,
		Port:     s.cfg.redisMasterPort,
		Password: s.cfg.redisMasterPassword,
		Db:       s.cfg.redisMasterDb,
	})
	if err != nil {
		return err
	}

	s.redisClient = db.Connect(ctx)

	return nil
}

func (s *Server) withQueueClient() error {
	queueClient, err := redisQueue.NewRedis(&redisQueue.RedisConfig{
		Host:     s.cfg.asynqRedisMasterHost,
		Port:     s.cfg.asynqRedisMasterPort,
		Password: s.cfg.asynqRedisMasterPassword,
		Db:       s.cfg.asynqRedisMasterDb,
	})
	if err != nil {
		return err
	}

	s.queueClient = queueClient.Create()

	return nil
}

func (s *Server) withWorker(ctx context.Context) error {
	var err error
	s.worker, err = worker.NewWorker(&worker.Config{
		AsynqRedisMasterHost:     s.cfg.asynqRedisMasterHost,
		AsynqRedisMasterPassword: s.cfg.asynqRedisMasterPassword,
		AsynqRedisMasterPort:     s.cfg.asynqRedisMasterPort,
		AsynqRedisMasterDb:       s.cfg.asynqRedisMasterDb,
	})
	if err != nil {
		return err
	}

	return nil
}

// func RunApp(cfg config) {
// 	// init database
// 	pgPool := initDatabase(cfg)
// 	defer pgPool.Close()

// 	// register to resolver
// 	pgResolver := postgres.NewResolver(pgPool)
// 	tenants, err := getAllTenant(pgPool)
// 	if err != nil {
// 		if !errors.Is(err, pgx.ErrNoRows) {
// 			log.Fatal("[x] Error during get all tenant")
// 		}
// 	}
// 	for _, tenant := range tenants {
// 		// store to resolver
// 		tenantPool, err := makeTenantConnection(tenant.Domain)
// 		if err != nil {
// 			log.Fatal("[x] error during creating tenant connection", err)
// 		}

// 		pgResolver.Register(postgres.Database{
// 			DomainName: postgres.Domain(tenant.Domain),
// 			Pool:       tenantPool,
// 		})
// 	}

// 	// init worker db conn
// 	workerDBConn := initWorkerDatabaseConnn(cfg)
// 	defer workerDBConn.Close()

// 	// init queue client
// 	queueClient := initQueueClient(cfg)
// 	defer func() {
// 		if err := queueClient.Close(); err != nil {
// 			log.Println("[x] error during closing queue connection")
// 		}
// 	}()

// 	//=== INITIALIZE REDIS ===
// 	rdb := initRedis(cfg)

// 	// create the http server
// 	app := NewApp(AppConfig{
// 		DB:               pgPool,
// 		Redis:            rdb,
// 		QueueClient:      queueClient,
// 		PostgresResolver: pgResolver,

// 		FiberCfg: fiber.Config{
// 			AppName: cfg.appName,
// 		},
// 	})

// 	// create worker server
// 	worker, err := NewWorker(workerDBConn, pgResolver)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	var sigChan = make(chan os.Signal, 1)
// 	signal.Notify(sigChan, os.Interrupt)

// 	// start the app
// 	var errChan = make(chan error, 1)
// 	go func() {
// 		err := app.Listen(cfg.httpPort)

// 		if err != nil {
// 			errChan <- err
// 		}
// 	}()

// 	go func() {
// 		log.Println("worker is running...")
// 		err := worker.Server.Run(worker.Mux)

// 		if err != nil {
// 			errChan <- err
// 		}
// 	}()

// 	select {
// 	case err := <-errChan:
// 		log.Fatal(err)
// 	case <-sigChan:
// 		var shutdownWg sync.WaitGroup

// 		go func() {
// 			shutdownWg.Add(1)
// 			defer shutdownWg.Done()
// 			// loop trhough db connection
// 			log.Println("closing database resolver connection")
// 			masterConn := pgResolver.MasterConn
// 			masterConn.Close()

// 			allTenantConn := pgResolver.GetAllTenantConn()
// 			if allTenantConn != nil {
// 				allTenantConn.Range(func(key, val interface{}) bool {
// 					log.Printf("closing %s database connection\n", key)
// 					pool := val.(*pgxpool.Pool)
// 					pool.Close()
// 					return true
// 				})
// 			}
// 		}()

// 		go func() {
// 			shutdownWg.Add(1)
// 			defer shutdownWg.Done()

// 			log.Println("shutdown app...")
// 			if err := app.Shutdown(); err != nil {
// 				log.Fatal("failed to shutdown", err)
// 			} else {
// 				log.Println("application is shutdown properly")
// 			}
// 		}()
// 		shutdownWg.Wait()
// 		log.Println("done!")
// 	}
// }

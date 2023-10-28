package tenant

import (
	"hris/module/auth/tenant/db"
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"
)

type Tenant struct {
	db *db.Db
}

type Dependency struct {
	PgResolver *postgres.Resolver
	Redis      *redisClient.Client
	MasterConn *pgxpool.Pool
}

func New(d *Dependency) *Tenant {
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth/tenant module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth/tenant module")
	}
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth/tenant module")
	}

	dbImpl := db.New(&db.Dependency{
		PgResolver: d.PgResolver,
		Redis:      d.Redis,
		MasterConn: d.MasterConn,
	})

	return &Tenant{
		db: dbImpl,
	}
}

package tenant

import (
	"hris/module/auth/tenant/db"
	"hris/module/shared/postgres"
	"log"

	redisClient "github.com/redis/go-redis/v9"
)

type Tenant struct {
	db *db.Db
}

type Dependency struct {
	PgResolver *postgres.Resolver
	Redis      *redisClient.Client
}

func New(d *Dependency) *Tenant {
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth/tenant module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth/tenant module")
	}

	return &Tenant{
		db: db.New(&db.Dependency{
			PgResolver: d.PgResolver,
			Redis:      d.Redis,
		}),
	}
}

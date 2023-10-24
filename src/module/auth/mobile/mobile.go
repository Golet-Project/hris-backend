package mobile

import (
	"hris/module/auth/mobile/db"
	"hris/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"

	userService "hris/module/user/service"
)

type Mobile struct {
	db *db.Db

	// other service
	userService *userService.Service
}

type Dependency struct {
	MasterConn *pgxpool.Pool
	PgResolver *postgres.Resolver
	Redis      *redisClient.Client

	// other service
	UserService *userService.Service
}

func New(d *Dependency) *Mobile {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth module")
	}
	if d.Redis == nil {
		log.Fatal("[x] Redis connection required on auth module")
	}

	return &Mobile{
		db: db.New(&db.Dependency{
			MasterConn: d.MasterConn,
			PgResolver: d.PgResolver,
			Redis:      d.Redis,
		}),

		userService: d.UserService,
	}
}

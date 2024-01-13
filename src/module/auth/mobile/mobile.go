package mobile

import (
	"hroost/module/auth/mobile/db"
	"hroost/module/auth/mobile/redis"
	"hroost/module/shared/postgres"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"

	userService "hroost/module/user/service"
)

type Mobile struct {
	db    *db.Db
	redis *redis.Redis

	// other service
	userService *userService.Service
}

type Dependency struct {
	MasterConn  *pgxpool.Pool
	PgResolver  *postgres.Resolver
	RedisClient *redisClient.Client

	// other service
	UserService *userService.Service
}

func New(d *Dependency) *Mobile {
	if d.MasterConn == nil {
		log.Fatal("[x] Master database connection required on auth/mobile module")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Database resolver required on auth/mobile module")
	}
	if d.RedisClient == nil {
		log.Fatal("[x] Redis connection required on auth/mobile module")
	}

	return &Mobile{
		db: db.New(&db.Dependency{
			MasterConn: d.MasterConn,
			PgResolver: d.PgResolver,
			Redis:      d.RedisClient,
		}),

		redis: redis.New(&redis.Dependency{
			Client: d.RedisClient,
		}),

		userService: d.UserService,
	}
}

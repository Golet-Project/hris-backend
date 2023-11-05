package auth

import (
	"log"

	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/auth/presentation/rest"
	"hris/module/auth/tenant"
	"hris/module/shared/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/redis/go-redis/v9"

	userService "hris/module/user/service"
)

type Auth struct {
	AuthPresentation *rest.AuthPresentation
}

type Dependency struct {
	DB          *pgxpool.Pool
	PgResolver  *postgres.Resolver
	RedisClient *redisClient.Client

	// other module service
	UserService *userService.Service
}

func InitAuth(d *Dependency) *Auth {
	if d.DB == nil {
		log.Fatal("[x] Auth package require a database connection")
	}
	if d.RedisClient == nil {
		log.Fatal("[x] Auth packge require a redis connection")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Auth package require a database resolver")
	}
	if d.UserService == nil {
		log.Fatal("[x] Auth package require a user service")
	}

	internalAuthService := internal.New(&internal.Dependency{
		Pg:    d.DB,
		Redis: d.RedisClient,
	})
	mobileAuthService := mobile.New(&mobile.Dependency{
		MasterConn: d.DB,
		PgResolver: d.PgResolver,

		RedisClient: d.RedisClient,

		UserService: d.UserService,
	})
	tenantAuthService := tenant.New(&tenant.Dependency{
		PgResolver: d.PgResolver,
		Redis:      d.RedisClient,
		MasterConn: d.DB,
	})

	authPresentation := rest.New(internalAuthService, mobileAuthService, tenantAuthService)

	return &Auth{
		AuthPresentation: authPresentation,
	}
}

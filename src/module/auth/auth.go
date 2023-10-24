package auth

import (
	"log"

	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/auth/presentation/rest"
	"hris/module/shared/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	userService "hris/module/user/service"
)

type Auth struct {
	AuthPresentation *rest.AuthPresentation
}

type Dependency struct {
	DB         *pgxpool.Pool
	PgResolver *postgres.Resolver
	Redis      *redis.Client

	// other module service
	userService *userService.Service
}

func InitAuth(d *Dependency) *Auth {
	if d.DB == nil {
		log.Fatal("[x] Auth package require a database connection")
	}
	if d.Redis == nil {
		log.Fatal("[x] Auth packge require a redis connection")
	}
	if d.PgResolver == nil {
		log.Fatal("[x] Auth package require a database resolver")
	}

	internalAuthService := internal.New(&internal.Dependency{
		Pg:    d.DB,
		Redis: d.Redis,
	})
	mobileAuthService := mobile.New(&mobile.Dependency{
		MasterConn: d.DB,
		PgResolver: d.PgResolver,
		Redis:      d.Redis,

		UserService: d.userService,
	})
	// webAuthService := service.NewWebAuthService(&authRepo)
	// mobileAuthService := service.NewMobileAuthService(&authRepo)

	return &Auth{
		AuthPresentation: &rest.AuthPresentation{
			Internal: internalAuthService,
			Mobile:   mobileAuthService,
		},
	}
}

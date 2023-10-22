package auth

import (
	"log"

	"hris/module/auth/internal"
	"hris/module/auth/presentation/rest"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Auth struct {
	AuthPresentation *rest.AuthPresentation
}

type Dependency struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func InitAuth(d *Dependency) *Auth {
	if d.DB == nil {
		log.Fatal("[x] Auth package require a database connection")
	}
	if d.Redis == nil {
		log.Fatal("[x] Auth packge require a redis connection")
	}

	internalAuthService := internal.New(&internal.Dependency{
		Pg:    d.DB,
		Redis: d.Redis,
	})
	// webAuthService := service.NewWebAuthService(&authRepo)
	// mobileAuthService := service.NewMobileAuthService(&authRepo)

	return &Auth{
		AuthPresentation: &rest.AuthPresentation{
			Internal: internalAuthService,
		},
	}
}

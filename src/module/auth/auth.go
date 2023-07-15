package auth

import (
	"hris/module/auth/presentation/rest"
	"hris/module/auth/repo/auth"
	"hris/module/auth/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Auth struct {
	AuthPresenter *rest.AuthPresenter
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

	authRepo := auth.Repository{
		DB:    d.DB,
		Redis: d.Redis,
	}

	internalAuthService := service.NewInternalAuthService(&authRepo)
	webAuthService := service.NewWebAuthService(&authRepo)
	mobileAuthService := service.NewMobileAuthService(&authRepo)

	return &Auth{
		AuthPresenter: &rest.AuthPresenter{
			InternalAuthService: internalAuthService,
			WebAuthService: webAuthService,
			MobileAuthService: mobileAuthService,
		},
	}
}

package auth

import (
	"hris/module/auth/presentation/rest"
	"hris/module/auth/repo/auth"
	"hris/module/auth/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	AuthPresenter *rest.AuthPresenter
}

type Dependency struct {
	DB *pgxpool.Pool
}

func InitAuth(d *Dependency) *Auth {
	authRepo := auth.Repository{
		DB: d.DB,
	}

	authService := service.NewAuthService(&authRepo)

	return &Auth{
		AuthPresenter: &rest.AuthPresenter{
			AuthService: authService,
		},
	}
}

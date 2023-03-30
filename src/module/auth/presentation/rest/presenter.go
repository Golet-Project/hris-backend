package rest

import (
	"hris/module/auth/service"
)

type AuthPresenter struct {
	AuthService *service.AuthService
}

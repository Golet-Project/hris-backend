package rest

import (
	"hris/module/auth/service"
)

type AuthPresenter struct {
	InternalAuthService *service.InternalAuthService
	WebAuthService *service.WebAuthService
	MobileAuthService *service.MobileAuthService
}

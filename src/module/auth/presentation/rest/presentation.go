package rest

import (
	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/auth/tenant"
)

type AuthPresentation struct {
	Internal *internal.Internal
	Mobile   *mobile.Mobile
	Tenant   *tenant.Tenant
}

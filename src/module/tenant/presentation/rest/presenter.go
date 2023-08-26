package rest

import "hris/module/tenant/service"

type TenantPresenter struct {
	Internal_TenantService *service.Internal_TenantService
}

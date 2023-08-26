package service

import "hris/module/tenant/repo/tenant"

type Internal_TenantService struct {
	TenantRepo *tenant.Repository
}

func NewInternal_TenantService(tenantRepo *tenant.Repository) *Internal_TenantService {
	return &Internal_TenantService{
		TenantRepo: tenantRepo,
	}
}

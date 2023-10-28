package rest

import (
	"hris/module/employee/tenant"
	"log"
)

type EmployeePresentation struct {
	tenant *tenant.Tenant
}

func New(tenant *tenant.Tenant) *EmployeePresentation {
	if tenant == nil {
		log.Fatal("[x] Tenant service required on employee presentation")
	}

	return &EmployeePresentation{
		tenant: tenant,
	}
}

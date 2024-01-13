package rest

import (
	"hroost/module/employee/mobile"
	"hroost/module/employee/tenant"
	"log"
)

type EmployeePresentation struct {
	tenant *tenant.Tenant
	mobile *mobile.Mobile
}

func New(tenant *tenant.Tenant, mobile *mobile.Mobile) *EmployeePresentation {
	if tenant == nil {
		log.Fatal("[x] Tenant service required on employee presentation")
	}
	if mobile == nil {
		log.Fatal("[x] Mobile service required on employee presentation")
	}

	return &EmployeePresentation{
		tenant: tenant,
		mobile: mobile,
	}
}

package rest

import (
	"hris/module/attendance/mobile"
	"hris/module/attendance/tenant"
	"log"
)

type AttandancePresentation struct {
	mobile *mobile.Mobile
	tenant *tenant.Tenant
}

func New(mobile *mobile.Mobile, tenant *tenant.Tenant) *AttandancePresentation {
	if mobile == nil {
		log.Fatal("[x] Mobile service required on attendance presentation")
	}
	if tenant == nil {
		log.Fatal("[x] Tenant service required on attendance presentation")
	}

	return &AttandancePresentation{
		mobile: mobile,
		tenant: tenant,
	}
}

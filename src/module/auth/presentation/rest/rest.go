package rest

import (
	"hris/module/auth/central"
	"hris/module/auth/mobile"
	"hris/module/auth/tenant"
	"log"
)

type AuthPresentation struct {
	central *central.Central
	mobile  *mobile.Mobile
	tenant  *tenant.Tenant
}

func New(central *central.Central, mobile *mobile.Mobile, tenant *tenant.Tenant) *AuthPresentation {
	if central == nil {
		log.Fatal("[x] Central service required on auth presentation")
	}
	if mobile == nil {
		log.Fatal("[x] Mobile service required on auth presentation")
	}
	if tenant == nil {
		log.Fatal("[x] Tenant service required on auth presentation")
	}

	return &AuthPresentation{
		central: central,
		mobile:  mobile,
		tenant:  tenant,
	}
}

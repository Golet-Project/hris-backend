package rest

import (
	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/auth/tenant"
	"log"
)

type AuthPresentation struct {
	internal *internal.Internal
	mobile   *mobile.Mobile
	tenant   *tenant.Tenant
}

func New(internal *internal.Internal, mobile *mobile.Mobile, tenant *tenant.Tenant) *AuthPresentation {
	if internal == nil {
		log.Fatal("[x] Internal service required on auth presentation")
	}
	if mobile == nil {
		log.Fatal("[x] Mobile service required on auth presentation")
	}
	if tenant == nil {
		log.Fatal("[x] Tenant service required on auth presentation")
	}

	return &AuthPresentation{
		internal: internal,
		mobile:   mobile,
		tenant:   tenant,
	}
}

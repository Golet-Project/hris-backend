package rest

import "hroost/module/tenant/central"

type TenantPresentation struct {
	central *central.Central
}

func New(central *central.Central) *TenantPresentation {
	if central == nil {
		panic("[x] Central service required on tenant presentation")
	}

	return &TenantPresentation{
		central: central,
	}
}

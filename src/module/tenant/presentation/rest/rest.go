package rest

import "hris/module/tenant/internal"

type TenantPresentation struct {
	internal *internal.Internal
}

func New(internal *internal.Internal) *TenantPresentation {
	if internal == nil {
		panic("[x] Internal service required on tenant presentation")
	}

	return &TenantPresentation{
		internal: internal,
	}
}

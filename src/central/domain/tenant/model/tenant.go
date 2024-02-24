package model

type CountTenantByDomainOut struct {
	Count int64
}

type CreateTenantIn struct {
	Name   string
	Domain string
}

type CreateTenantOut struct {
	UID    string
	Name   string
	Domain string
}

type MigrateTenantDBIn struct {
	Domain string `json:"domain"`
}

package tenant

import (
	"hris/module/tenant/presentation/rest"
	"hris/module/tenant/repo/tenant"
	"hris/module/tenant/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	Internal_TenantPresenter *rest.TenantPresenter
}

type Dependency struct {
	DB *pgxpool.Pool
}

func InitTenant(d *Dependency) *Tenant {
	if d.DB == nil {
		log.Fatal("[x] Database connection required on tenant module")
	}

	// init repo
	internal_tenantRepo := tenant.Repository{
		DB: d.DB,
	}

	// init service
	internal_tenantService := service.NewInternal_TenantService(&internal_tenantRepo)

	return &Tenant{
		Internal_TenantPresenter: &rest.TenantPresenter{
			Internal_TenantService: internal_tenantService,
		},
	}
}

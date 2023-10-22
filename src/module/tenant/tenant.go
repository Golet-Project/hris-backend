package tenant

import (
	"hris/module/tenant/internal"
	"hris/module/tenant/presentation/rest"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	TenantPresentation *rest.TenantPresentation
}

type Dependency struct {
	DB          *pgxpool.Pool
	QueueClient *asynq.Client
}

func InitTenant(d *Dependency) *Tenant {
	if d.DB == nil {
		log.Fatal("[x] Database connection required on tenant module")
	}
	if d.QueueClient == nil {
		log.Fatal("[x] Queue client required on tenant module")
	}

	// init service
	internalTenantService := internal.New(&internal.Dependency{
		Pg:    d.DB,
		Queue: d.QueueClient,
	})

	return &Tenant{
		TenantPresentation: &rest.TenantPresentation{
			Internal: internalTenantService,
		},
	}
}

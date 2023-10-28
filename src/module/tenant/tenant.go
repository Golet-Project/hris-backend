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
	MasterConn  *pgxpool.Pool
	QueueClient *asynq.Client
}

func InitTenant(d *Dependency) *Tenant {
	if d.MasterConn == nil {
		log.Fatal("[x] Database connection required on tenant module")
	}
	if d.QueueClient == nil {
		log.Fatal("[x] Queue client required on tenant module")
	}

	// init service
	internalTenantService := internal.New(&internal.Dependency{
		MasterConn: d.MasterConn,
		Queue:      d.QueueClient,
	})

	tenantPresentation := rest.New(internalTenantService)

	return &Tenant{
		TenantPresentation: tenantPresentation,
	}
}

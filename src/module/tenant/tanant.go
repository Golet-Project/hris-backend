package tenant

import (
	"hris/module/tenant/presentation/rest"
	"hris/module/tenant/repo/queue"
	"hris/module/tenant/repo/tenant"
	"hris/module/tenant/service"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tenant struct {
	Internal_TenantPresenter *rest.TenantPresenter
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

	// init repo
	internal_tenantRepo := tenant.Repository{
		DB: d.DB,
	}
	queueRepo := queue.QueueRepo{
		Client: d.QueueClient,
	}

	// init service
	internal_tenantService := service.NewInternal_TenantService(&internal_tenantRepo, &queueRepo)

	return &Tenant{
		Internal_TenantPresenter: &rest.TenantPresenter{
			Internal_TenantService: internal_tenantService,
		},
	}
}

package service

import (
	"hris/module/tenant/repo/queue"
	"hris/module/tenant/repo/tenant"
)

type Internal_TenantService struct {
	TenantRepo *tenant.Repository

	QueueRepo *queue.QueueRepo
}

func NewInternal_TenantService(tenantRepo *tenant.Repository, queueRepo *queue.QueueRepo) *Internal_TenantService {
	return &Internal_TenantService{
		TenantRepo: tenantRepo,
		QueueRepo: queueRepo,
	}
}

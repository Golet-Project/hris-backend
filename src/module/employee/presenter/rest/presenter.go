package rest

import "hris/module/employee/tenant"

type EmployeePresentation struct {
	Tenant *tenant.Tenant
}

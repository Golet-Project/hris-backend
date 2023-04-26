package rest

import "hris/module/web/employee/service"

type EmployeePresenter struct {
	EmployeeService *service.EmployeeService
}

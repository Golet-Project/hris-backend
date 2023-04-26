package service

import "hris/module/web/employee/repo/employee"

type EmployeeService struct {
	EmployeeRepo *employee.Repository
}

func NewEmployeeService(EmployeeRepo *employee.Repository) *EmployeeService {
	return &EmployeeService{
		EmployeeRepo: EmployeeRepo,
	}
}

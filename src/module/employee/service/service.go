package service

import "hris/module/employee/repo/employee"

type WebEmployeeService struct {
	EmployeeRepo *employee.Repository
}

func NewWebEmployeeService(employeeRepo *employee.Repository) *WebEmployeeService {
	return &WebEmployeeService{
		EmployeeRepo: employeeRepo,
	}
}
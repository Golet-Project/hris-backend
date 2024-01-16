package server

import (
	"hroost/presentation/rest/attendance"
	"hroost/presentation/rest/auth"
	"hroost/presentation/rest/employee"
	"hroost/presentation/rest/homepage"
	"hroost/presentation/rest/region"
	"hroost/presentation/rest/tenant_management"
)

type MobileServiceProvider struct {
}

type AppProvider struct {
	attendanceRest       *attendance.Attendance
	authRest             *auth.Auth
	employeeRest         *employee.Employee
	homepageRest         *homepage.Homepage
	regionRest           *region.Region
	tenantManagementRest *tenant_management.TenantManagement
}

func (s *Server) newAppProvider() *AppProvider {

	attendanceRest := attendance.NewAttendance(&attendance.Config{})
}

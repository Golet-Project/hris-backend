package attendance

import (
	"hris/module/attendance/mobile"
	"hris/module/attendance/presentation/rest"
	"hris/module/attendance/tenant"
	"hris/module/shared/postgres"
	"log"

	userService "hris/module/user/service"
)

type Attendance struct {
	AttendancePresentation *rest.AttandancePresentation
}

type Dependency struct {
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func InitAtteandance(d *Dependency) *Attendance {
	if d.PgResolver == nil {
		log.Fatal("[x] postgres resolver is required on attendance module")
	}
	if d.UserService == nil {
		log.Fatal("[x] user service is required on attendance module")
	}

	mobileAttendanceService := mobile.New(&mobile.Dependency{
		PgResolver: d.PgResolver,

		UserService: d.UserService,
	})
	tenantAttendanceService := tenant.New(&tenant.Dependency{
		PgResolver: d.PgResolver,

		UserService: d.UserService,
	})

	attendancePresentation := rest.New(mobileAttendanceService, tenantAttendanceService)

	return &Attendance{
		AttendancePresentation: attendancePresentation,
	}
}

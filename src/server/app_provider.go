package server

import (
	"hroost/presentation/rest/attendance"
	"hroost/presentation/rest/auth"
	"hroost/presentation/rest/employee"
	"hroost/presentation/rest/homepage"
	"hroost/presentation/rest/region"
	"hroost/presentation/rest/tenant_management"

	sharedRegionDbProvince "hroost/shared/domain/region/db/province"
	sharedUserDb "hroost/shared/domain/user/db"

	sharedRegionService "hroost/shared/domain/region/service"
	sharedUserService "hroost/shared/domain/user/service"

	centralAuthDb "hroost/central/domain/auth/db"
	centralAuthMemory "hroost/central/domain/auth/memory"
	centralTenantDb "hroost/central/domain/tenant/db"
	centralTenantQueue "hroost/central/domain/tenant/queue"

	centralAuthService "hroost/central/domain/auth/service"
	centralTenantService "hroost/central/domain/tenant/service"

	mobileAttendanceDb "hroost/mobile/domain/attendance/db"
	mobileAuthDb "hroost/mobile/domain/auth/db"
	mobileAuthMemory "hroost/mobile/domain/auth/memory"
	mobileEmployeeDb "hroost/mobile/domain/employee/db"
	mobileHomepageDb "hroost/mobile/domain/homepage/db"

	mobileAttendanceService "hroost/mobile/domain/attendance/service"
	mobileAuthService "hroost/mobile/domain/auth/service"
	mobileEmployeeService "hroost/mobile/domain/employee/service"
	mobileHomepageService "hroost/mobile/domain/homepage/service"

	tenantAttendanceDb "hroost/tenant/domain/attendance/db"
	tenantAuthDb "hroost/tenant/domain/auth/db"
	tenantEmployeeDb "hroost/tenant/domain/employee/db"

	tenantAttendanceService "hroost/tenant/domain/attendance/service"
	tenantAuthService "hroost/tenant/domain/auth/service"
	tenantEmployeeService "hroost/tenant/domain/employee/service"
)

type SharedServiceProvider struct {
	regionService *sharedRegionService.Service
	userService   *sharedUserService.Service
}

func (s *Server) initShared() (*SharedServiceProvider, error) {
	// region
	regionDbProvince, err := sharedRegionDbProvince.New(&sharedRegionDbProvince.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	regionService, err := sharedRegionService.New(&sharedRegionService.Config{ProvinceDb: regionDbProvince})
	if err != nil {
		return nil, err
	}

	// user
	userDb, err := sharedUserDb.New(&sharedUserDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}
	userService, err := sharedUserService.New(&sharedUserService.Config{Db: userDb})
	if err != nil {
		return nil, err
	}

	return &SharedServiceProvider{
		regionService: regionService,
		userService:   userService,
	}, nil
}

type CentralServiceProvider struct {
	authService   *centralAuthService.Service
	tenantSerivce *centralTenantService.Service
}

func (s *Server) initCentral() (*CentralServiceProvider, error) {
	// auth
	authDb, err := centralAuthDb.New(&centralAuthDb.Config{
		Redis:      s.redisClient,
		PgResolver: s.pgResolver,
	})
	if err != nil {
		return nil, err
	}
	authMemory, err := centralAuthMemory.New(&centralAuthMemory.Config{Client: s.redisClient})
	if err != nil {
		return nil, err
	}
	authService, err := centralAuthService.New(&centralAuthService.Config{
		Db:     authDb,
		Memory: authMemory,
	})

	// tenant
	tenantDb, err := centralTenantDb.New(&centralTenantDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	tenantQueue, err := centralTenantQueue.New(&centralTenantQueue.Config{Client: s.queueClient})
	if err != nil {
		return nil, err
	}
	tenantService, err := centralTenantService.New(&centralTenantService.Config{
		Db:    tenantDb,
		Queue: tenantQueue,
	})
	if err != nil {
		return nil, err
	}

	return &CentralServiceProvider{
		authService:   authService,
		tenantSerivce: tenantService,
	}, nil
}

type MobileServiceProvider struct {
	attendanceService *mobileAttendanceService.Service
	authService       *mobileAuthService.Service
	employeeService   *mobileEmployeeService.Service
	homepageService   *mobileHomepageService.Service
}

func (s *Server) initMobile() (*MobileServiceProvider, error) {
	sharedService, err := s.initShared()
	if err != nil {
		return nil, err
	}

	// attendance
	attendanceDb, err := mobileAttendanceDb.New(&mobileAttendanceDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	attendanceService, err := mobileAttendanceService.New(&mobileAttendanceService.Config{
		Db: attendanceDb,

		// shared service
		UserService: sharedService.userService,
	})
	if err != nil {
		return nil, err
	}

	// auth
	authDb, err := mobileAuthDb.New(&mobileAuthDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}
	authMemory, err := mobileAuthMemory.New(&mobileAuthMemory.Config{Client: s.redisClient})
	if err != nil {
		return nil, err
	}
	authService, err := mobileAuthService.New(&mobileAuthService.Config{
		Db:     authDb,
		Memory: authMemory,

		// shared service
		UserService: sharedService.userService,
	})

	// employee
	employeeDb, err := mobileEmployeeDb.New(&mobileEmployeeDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	employeeService, err := mobileEmployeeService.New(&mobileEmployeeService.Config{
		Db: employeeDb,

		// shared service
		UserService: sharedService.userService,
	})
	if err != nil {
		return nil, err
	}

	// homepage
	homepageDb, err := mobileHomepageDb.New(&mobileHomepageDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	homepageService, err := mobileHomepageService.New(&mobileHomepageService.Config{
		Db: homepageDb,

		// shared service
		UserService: sharedService.userService,
	})

	return &MobileServiceProvider{
		attendanceService: attendanceService,
		authService:       authService,
		employeeService:   employeeService,
		homepageService:   homepageService,
	}, nil
}

type TenantServiceProvider struct {
	attendanceService *tenantAttendanceService.Service
	authService       *tenantAuthService.Service
	employeeService   *tenantEmployeeService.Service
}

func (s *Server) initTenant() (*TenantServiceProvider, error) {
	sharedService, err := s.initShared()
	if err != nil {
		return nil, err
	}

	// attendance
	attendanceDb, err := tenantAttendanceDb.New(&tenantAttendanceDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	attendanceService, err := tenantAttendanceService.New(&tenantAttendanceService.Config{
		Db: attendanceDb,

		// shared service
		UserService: sharedService.userService,
	})
	if err != nil {
		return nil, err
	}

	// auth
	authDb, err := tenantAuthDb.New(&tenantAuthDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}
	authService, err := tenantAuthService.New(&tenantAuthService.Config{Db: authDb})
	if err != nil {
		return nil, err
	}

	// employee
	employeeDb, err := tenantEmployeeDb.New(&tenantEmployeeDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	employeeService, err := tenantEmployeeService.New(&tenantEmployeeService.Config{Db: employeeDb})
	if err != nil {
		return nil, err
	}

	return &TenantServiceProvider{
		attendanceService: attendanceService,
		authService:       authService,
		employeeService:   employeeService,
	}, nil
}

type restPresentation struct {
	attendance       *attendance.Attendance
	auth             *auth.Auth
	employee         *employee.Employee
	homepage         *homepage.Homepage
	region           *region.Region
	tenantManagement *tenant_management.TenantManagement
}

type Presentation struct {
	rest *restPresentation
}

func (s *Server) newAppProvider() (*Presentation, error) {
	sharedServiceProvider, err := s.initShared()
	if err != nil {
		return nil, err
	}

	centralServiceProvider, err := s.initCentral()
	if err != nil {
		return nil, err
	}

	mobileServiceProvider, err := s.initMobile()
	if err != nil {
		return nil, err
	}

	tenantServiceProvider, err := s.initTenant()
	if err != nil {
		return nil, err
	}

	// attendance
	attendanceRest, err := attendance.NewAttendance(&attendance.Config{
		MobileService: mobileServiceProvider.attendanceService,
		TenantService: tenantServiceProvider.attendanceService,
	})
	if err != nil {
		return nil, err
	}

	// auth
	authRest, err := auth.NewAuth(&auth.Config{
		CentralService: centralServiceProvider.authService,
		MobileService:  mobileServiceProvider.authService,
		TenantService:  tenantServiceProvider.authService,
	})
	if err != nil {
		return nil, err
	}

	// employee
	employeeRest, err := employee.NewEmployee(&employee.Config{
		TenantService: tenantServiceProvider.employeeService,
		MobileService: mobileServiceProvider.employeeService,
	})
	if err != nil {
		return nil, err
	}

	// homepage
	homepageRest, err := homepage.NewHomepage(&homepage.Config{
		MobileService: mobileServiceProvider.homepageService,
	})
	if err != nil {
		return nil, err
	}

	// region
	regionRest, err := region.NewRegion(&region.Config{
		Service: sharedServiceProvider.regionService,
	})

	// tenantManagement
	tenantManagementRest, err := tenant_management.NewTenantManagement(&tenant_management.Config{
		CentralService: centralServiceProvider.tenantSerivce,
	})

	return &Presentation{
		rest: &restPresentation{
			attendance:       attendanceRest,
			auth:             authRest,
			employee:         employeeRest,
			homepage:         homepageRest,
			region:           regionRest,
			tenantManagement: tenantManagementRest,
		},
	}, nil
}

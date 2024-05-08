package server

import (
	"hroost/presentation/rest/attendance"
	"hroost/presentation/rest/auth"
	"hroost/presentation/rest/employee"
	"hroost/presentation/rest/homepage"
	"hroost/presentation/rest/region"
	"hroost/presentation/rest/tenant_management"
	"os"

	sharedRegionDbDistrict "hroost/shared/domain/region/db/district"
	sharedRegionDbProvince "hroost/shared/domain/region/db/province"
	sharedRegionDbRegency "hroost/shared/domain/region/db/regency"
	sharedRegionDbVillage "hroost/shared/domain/region/db/village"
	sharedUserDb "hroost/shared/domain/user/db"

	sharedRegionService "hroost/shared/domain/region/service"
	sharedUserService "hroost/shared/domain/user/service"

	centralAuthDb "hroost/central/domain/auth/db"
	centralAuthMemory "hroost/central/domain/auth/memory"
	centralTenantDb "hroost/central/domain/tenant/db"
	centralTenantQueue "hroost/central/domain/tenant/queue"

	mobileAttendanceDb "hroost/mobile/domain/attendance/db"
	mobileAuthDb "hroost/mobile/domain/auth/db"
	mobileAuthMemory "hroost/mobile/domain/auth/memory"
	mobileEmployeeDb "hroost/mobile/domain/employee/db"
	mobileHomepageDb "hroost/mobile/domain/homepage/db"

	tenantAttendanceDb "hroost/tenant/domain/attendance/db"
	tenantAuthDb "hroost/tenant/domain/auth/db"
	tenantEmployeeDb "hroost/tenant/domain/employee/db"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
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
	regionDbRegency, err := sharedRegionDbRegency.New(&sharedRegionDbRegency.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	regionDbDistrict, err := sharedRegionDbDistrict.New(&sharedRegionDbDistrict.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	regionDbVillage, err := sharedRegionDbVillage.New(&sharedRegionDbVillage.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}

	regionService, err := sharedRegionService.New(&sharedRegionService.Config{
		ProvinceDb: regionDbProvince,
		RegencyDb:  regionDbRegency,
		DistrictDb: regionDbDistrict,
		VillageDb:  regionDbVillage,
	})
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
	authDb     *centralAuthDb.Db
	authMemory *centralAuthMemory.Memory

	tenantDb    *centralTenantDb.Db
	tenantQueue *centralTenantQueue.Queue

	oauth2Cfg *oauth2.Config
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

	// oauthcfg
	oauth2cfg := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Endpoint:     endpoints.Google,
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	// tenant
	tenantDb, err := centralTenantDb.New(&centralTenantDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}
	tenantQueue, err := centralTenantQueue.New(&centralTenantQueue.Config{Client: s.queueClient})
	if err != nil {
		return nil, err
	}

	return &CentralServiceProvider{
		authDb:     authDb,
		authMemory: authMemory,

		tenantDb:    tenantDb,
		tenantQueue: tenantQueue,

		oauth2Cfg: oauth2cfg,
	}, nil
}

type MobileServiceProvider struct {
	authDb     *mobileAuthDb.Db
	authMemory *mobileAuthMemory.Memory

	attendanceDb *mobileAttendanceDb.Db

	employeeDb *mobileEmployeeDb.Db

	homepageDb *mobileHomepageDb.Db
}

func (s *Server) initMobile() (*MobileServiceProvider, error) {
	// attendance
	attendanceDb, err := mobileAttendanceDb.New(&mobileAttendanceDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
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

	// employee
	employeeDb, err := mobileEmployeeDb.New(&mobileEmployeeDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}

	// homepage
	homepageDb, err := mobileHomepageDb.New(&mobileHomepageDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}

	return &MobileServiceProvider{
		authDb:     authDb,
		authMemory: authMemory,

		attendanceDb: attendanceDb,

		employeeDb: employeeDb,

		homepageDb: homepageDb,
	}, nil
}

type TenantServiceProvider struct {
	attendanceDb *tenantAttendanceDb.Db

	authDb *tenantAuthDb.Db

	employeeDb *tenantEmployeeDb.Db
}

func (s *Server) initTenant() (*TenantServiceProvider, error) {
	// attendance
	attendanceDb, err := tenantAttendanceDb.New(&tenantAttendanceDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}

	// auth
	authDb, err := tenantAuthDb.New(&tenantAuthDb.Config{PgResolver: s.pgResolver, Redis: s.redisClient})
	if err != nil {
		return nil, err
	}

	// employee
	employeeDb, err := tenantEmployeeDb.New(&tenantEmployeeDb.Config{PgResolver: s.pgResolver})
	if err != nil {
		return nil, err
	}

	return &TenantServiceProvider{
		attendanceDb: attendanceDb,

		authDb: authDb,

		employeeDb: employeeDb,
	}, nil
}

type restPresentation struct {
	attendance       *attendance.Attendance
	auth             *auth.Auth
	employee         *employee.Employee
	homepage         *homepage.HomePage
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
		Mobile: &attendance.Mobile{
			Db: mobileServiceProvider.attendanceDb,
		},

		Tenant: &attendance.Tenant{
			Db: tenantServiceProvider.attendanceDb,
		},
	})
	if err != nil {
		return nil, err
	}

	// auth
	authRest, err := auth.NewAuth(&auth.Config{
		Central: &auth.Central{
			Db:     centralServiceProvider.authDb,
			Memory: centralServiceProvider.authMemory,
		},

		Mobile: &auth.Mobile{
			Db:     mobileServiceProvider.authDb,
			Memory: mobileServiceProvider.authMemory,
		},

		Tenant: &auth.Tenant{
			Db: tenantServiceProvider.authDb,
		},
	})
	if err != nil {
		return nil, err
	}

	// employee
	employeeRest, err := employee.NewEmployee(&employee.Config{
		Mobile: &employee.Mobile{
			Db: mobileServiceProvider.employeeDb,
		},

		Tenant: &employee.Tenant{
			Db: tenantServiceProvider.employeeDb,
		},
	})
	if err != nil {
		return nil, err
	}

	// homepage
	homepageRest, err := homepage.NewHomepage(&homepage.Config{
		Mobile: &homepage.Mobile{
			Db: mobileServiceProvider.homepageDb,
		},
	})
	if err != nil {
		return nil, err
	}

	// region
	regionRest, err := region.NewRegion(&region.Config{
		Service: sharedServiceProvider.regionService,
	})
	if err != nil {
		return nil, err
	}

	// tenantManagement
	tenantManagementRest, err := tenant_management.NewTenantManagement(&tenant_management.Config{
		Central: &tenant_management.Central{
			Db:    centralServiceProvider.tenantDb,
			Queue: centralServiceProvider.tenantQueue,
		},
	})
	if err != nil {
		return nil, err
	}

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

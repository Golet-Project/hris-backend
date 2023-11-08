package homepage

import (
	"hris/module/homepage/mobile"
	"hris/module/homepage/presentation/rest"
	"hris/module/shared/postgres"
	userService "hris/module/user/service"
	"log"
)

type HomePage struct {
	Rest *rest.Rest
}

type Dependency struct {
	PgResolver *postgres.Resolver

	// other service
	UserService *userService.Service
}

func InitHomePage(d *Dependency) *HomePage {
	if d.PgResolver == nil {
		log.Fatal("[x] HomePage package require a database resolver")
	}
	if d.UserService == nil {
		log.Fatal("[x] HomePage package require a user service")
	}

	// tenantHomePageService := tenant.New(&tenant.Dependency{
	// 	MasterConn: d.MasterDB,
	// })
	mobileHomePageService := mobile.New(&mobile.Dependency{
		PgResolver: d.PgResolver,

		UserService: d.UserService,
	})

	// homePagePresentation := rest.New(tenantHomePageService, mobileHomePageService)

	return &HomePage{
		Rest: rest.New(mobileHomePageService),
	}
}

package cmd

import (
	business "hris/module/mobile/user/service"
	"hris/presentation/rest/mobile"
	"hris/presentation/rest/mobile/user"
	"hris/presentation/rest/web"

	"github.com/gofiber/fiber/v2"
)

// NewHttpServer create the http server
func NewHttpServer(cfg ...fiber.Config) *fiber.App {
	var app *fiber.App

	if len(cfg) > 0 {
		app = fiber.New(cfg[0])
	} else {
		app = fiber.New()
	}

	// service
	userService := business.UserService{}

	// intialize route
	mobile := mobile.Dependency{
		UserPresenter: &user.UserPresenter{
			UserService: &userService,
		},
	}

	web := web.Dependency{}

	mobile.Route(app.Group("/"))
	web.Route(app.Group("/w"))

	return app
}

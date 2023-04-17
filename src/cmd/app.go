package cmd

import (
	auth "hris/module/auth"
	authPresentation "hris/module/auth/presentation/rest"
	authRepo "hris/module/auth/repo/auth"
	authService "hris/module/auth/service"
	"reflect"

	"hris/module/mobile"
	"hris/module/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppConfig struct {
	DB       *pgxpool.Pool
	FiberCfg fiber.Config
}

// NewApp initialize the app
func NewApp(config AppConfig) *fiber.App {
	var app *fiber.App

	if !reflect.DeepEqual(config.FiberCfg, fiber.Config{}) {
		app = fiber.New(config.FiberCfg)
	} else {
		app = fiber.New()
	}

	app.Use(logger.New())
	//===== AUTH =====
	// initialize auth repo
	authRepo := &authRepo.Repository{
		DB: config.DB,
	}

	// initialize auth service
	authService := authService.NewAuthService(authRepo)

	// initialize auth route
	auth := auth.Dependency{
		AuthPresenter: &authPresentation.AuthPresenter{
			AuthService: authService,
		},
	}

	// initialize web route
	web := web.Dependency{}

	// intialize mobile route
	mobile := mobile.Dependency{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("oke")
	})

	auth.Route(app)
	mobile.Route(app.Group("/m"))
	web.Route(app.Group("/w"))

	return app
}

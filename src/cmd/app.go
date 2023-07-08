package cmd

import (
	auth "hris/module/auth"
	"reflect"

	employeeWeb "hris/module/web/employee"

	"hris/module/mobile"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	auth := auth.InitAuth(&auth.Dependency{
		DB: config.DB,
	})

	// initialize web's modules
	//===== EMPLOYEE =====
	employeeWeb := employeeWeb.InitEmployee(&employeeWeb.Dependency{
		DB: config.DB,
	})

	// intialize mobile route
	mobile := mobile.Dependency{}

	app.Use(cors.New(cors.Config{
		// AllowOrigins: "*",
		// AllowOrigins: "http://localhost:3000,https://google.com",
		// AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("oke")
	})

	auth.Route(app)
	mobile.Route(app.Group("/m"))
	employeeWeb.Route(app.Group("/w"))

	return app
}

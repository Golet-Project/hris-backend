package main

import (
	auth "hris/module/auth"
	"hris/module/employee"
	"hris/module/region"
	"hris/module/shared/primitive"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mileusna/useragent"
	"github.com/redis/go-redis/v9"
)

type AppConfig struct {
	DB       *pgxpool.Pool
	Redis    *redis.Client
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
		DB:    config.DB,
		Redis: config.Redis,
	})

	//=== EMPLOYEE ===
	employee := employee.InitEmployee(&employee.Dependency{
		DB: config.DB,
	})

	//=== Region ===
	region := region.InitRegion(&region.Dependency{
		DB: config.DB,
	})

	app.Use(cors.New(cors.Config{
		// AllowOrigins: "*",
		// AllowOrigins: "http://localhost:3000,https://google.com",
		// AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// check where the request is coming from, then translate it into an application ID
	app.Use(func(c *fiber.Ctx) error {
		userAgent := string(c.Context().UserAgent())

		ua := useragent.Parse(userAgent)

		switch ua.Name { // from browser
		case useragent.Opera:
			fallthrough
		case useragent.OperaMini:
			fallthrough
		case useragent.OperaTouch:
			fallthrough
		case useragent.Chrome:
			fallthrough
		case useragent.HeadlessChrome:
			fallthrough
		case useragent.Firefox:
			fallthrough
		case useragent.InternetExplorer:
			fallthrough
		case useragent.Safari:
			fallthrough
		case useragent.Edge:
			fallthrough
		case useragent.Vivaldi:
			appId := c.Get("X-App-ID")
			switch appId {
			case primitive.InternalAppID.String():
				c.Locals("AppID", primitive.InternalAppID)
				return c.Next()
			case primitive.WebAppID.String():
				c.Locals("AppID", primitive.WebAppID)
				return c.Next()

			default:
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(map[string]string{
					"message": "invalid app ID",
				})
			}

		default: // verify from mobile devices
			switch ua.OS {
			case useragent.Android:
				fallthrough
			case useragent.IOS:
				c.Locals("AppID", primitive.MobileAppID)
				return c.Next()

			default:
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(map[string]string{
					"message": "invalid client",
				})
			}
		}
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("oke")
	})

	auth.Route(app)
	employee.Route(app)
	region.Route(app)

	return app
}

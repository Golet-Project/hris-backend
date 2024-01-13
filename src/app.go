package main

import (
	"hroost/module/attendance"
	"hroost/module/auth"
	"hroost/module/employee"
	"hroost/module/homepage"
	"hroost/module/region"
	"hroost/module/shared/postgres"
	"hroost/module/shared/primitive"
	"hroost/module/tenant"
	"hroost/module/user"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mileusna/useragent"
	"github.com/redis/go-redis/v9"
)

type AppConfig struct {
	DB               *pgxpool.Pool
	Redis            *redis.Client
	FiberCfg         fiber.Config
	QueueClient      *asynq.Client
	PostgresResolver *postgres.Resolver
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
	//=== User ===
	user := user.InitUser(&user.Dependency{
		MasterDB:    config.DB,
		RedisClient: config.Redis,
	})

	//===== AUTH =====
	auth := auth.InitAuth(&auth.Dependency{
		PgResolver:  config.PostgresResolver,
		DB:          config.DB,
		RedisClient: config.Redis,

		UserService: user.UserService,
	})

	//=== EMPLOYEE ===
	employee := employee.InitEmployee(&employee.Dependency{
		MasterDB:   config.DB,
		PgResolver: config.PostgresResolver,

		UserService: user.UserService,
	})

	//=== Region ===
	region := region.InitRegion(&region.Dependency{
		DB: config.DB,
	})

	//=== Tenant ===
	tenant := tenant.InitTenant(&tenant.Dependency{
		MasterConn:  config.DB,
		QueueClient: config.QueueClient,
	})

	//=== Attendance ===
	attendance := attendance.InitAtteandance(&attendance.Dependency{
		MasterConn: config.DB,
		PgResolver: config.PostgresResolver,

		UserService: user.UserService,
	})

	//== Homepae ===
	homepage := homepage.InitHomePage(&homepage.Dependency{
		PgResolver: config.PostgresResolver,

		UserService: user.UserService,
	})

	app.Use(cors.New(cors.Config{
		// AllowOrigins: "*",
		// AllowOrigins: "http://localhost:3000,https://google.com",
		// AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// check where the request is coming from, then translate it into an application ID
	app.Use(func(c *fiber.Ctx) error {
		// except for the change password endpoint
		originalUrl := utils.CopyString(c.OriginalURL())
		if strings.HasPrefix(originalUrl, "/auth/password") && c.Method() == "PUT" {
			return c.Next()
		}

		userAgent := utils.CopyString(string(c.Context().UserAgent()))

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
			appId := utils.CopyString(c.Get("X-App-ID"))
			domain := utils.CopyString(c.Get("X-Domain"))
			switch appId {
			case primitive.CentralAppID.String():
				c.Locals("AppID", primitive.CentralAppID)
				return c.Next()
			case primitive.TenantAppID.String():
				c.Locals("AppID", primitive.TenantAppID)
				c.Locals("domain", domain)
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
	tenant.Route(app)
	attendance.Route(app)
	homepage.Route(app)

	return app
}

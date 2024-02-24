package employee

import (
	"fmt"

	"hroost/shared/primitive"

	mobileDb "hroost/mobile/domain/employee/db"
	mobileService "hroost/mobile/domain/employee/service"
	mobileJwt "hroost/mobile/lib/jwt"

	tenantDb "hroost/tenant/domain/employee/db"
	tenantService "hroost/tenant/domain/employee/service"
	tenantJwt "hroost/tenant/lib/jwt"

	"github.com/gofiber/fiber/v2"
)

type Mobile struct {
	Db *mobileDb.Db
}

type Tenant struct {
	Db *tenantDb.Db
}

type Config struct {
	Mobile *Mobile
	Tenant *Tenant
}

type Employee struct {
	mobile *Mobile
	tenant *Tenant
}

func NewEmployee(cfg *Config) (*Employee, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependecy for employee required")
	}
	if cfg.Mobile == nil {
		return nil, fmt.Errorf("mobile module required")
	}
	if cfg.Tenant == nil {
		return nil, fmt.Errorf("tenant module required")
	}

	return &Employee{
		mobile: cfg.Mobile,
		tenant: cfg.Tenant,
	}, nil
}

func (e Employee) CreateEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(tenantJwt.CustomClaims)

		var body tenantService.CreateEmployeeIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Domain = claims.Domain

		service := tenantService.CreateEmployee{
			Db: e.tenant.Db,
		}

		// call the services
		serviceOut := service.Exec(c.Context(), body)

		res.Message = serviceOut.GetMessage()
		res.Data = nil

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

func (e Employee) FindAllEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(tenantJwt.CustomClaims)

		var req tenantService.FindAllEmployeeIn
		req.Domain = claims.Domain

		service := tenantService.FindAllEmployee{
			Db: e.tenant.Db,
		}

		// call the services
		serviceOut := service.Exec(c.Context(), req)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 {
			res.Data = nil
		} else {
			data := make([]interface{}, len(serviceOut.Employee))
			for i, v := range serviceOut.Employee {
				data[i] = v
			}
			res.Data = data
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

func (e Employee) GetProfile(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.MobileAppID:
		claims := c.Locals("user_auth").(mobileJwt.CustomClaims)

		service := mobileService.GetProfile{
			Db: e.mobile.Db,
		}

		// call the service
		serviceOut := service.Exec(c.Context(), mobileService.GetProfileIn{
			UID: claims.UserUID,
		})

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() <= 500 {
			res.Data = nil
		} else {
			res.Data = serviceOut
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

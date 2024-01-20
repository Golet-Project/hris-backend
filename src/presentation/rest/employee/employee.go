package employee

import (
	"fmt"
	mobileJwt "hroost/mobile/lib/jwt"
	tenantJwt "hroost/tenant/lib/jwt"

	"hroost/shared/primitive"

	mobileService "hroost/mobile/domain/employee/service"
	tenantService "hroost/tenant/domain/employee/service"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	TenantService *tenantService.Service
	MobileService *mobileService.Service
}

type Employee struct {
	tenantService *tenantService.Service
	mobileService *mobileService.Service
}

func NewEmployee(cfg *Config) (*Employee, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependecy for employee required")
	}

	return &Employee{}, nil
}

func (e Employee) CreateEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(tenantJwt.CustomClaims)

		var req tenantService.FindAllEmployeeIn
		req.Domain = claims.Domain

		// call the services
		serviceOut := e.tenantService.FindAllEmployee(c.Context(), req)

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

func (e Employee) FindAllEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(tenantJwt.CustomClaims)

		var req tenantService.FindAllEmployeeIn
		req.Domain = claims.Domain

		// call the services
		serviceOut := e.tenantService.FindAllEmployee(c.Context(), req)

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

		// call the service
		serviceOut := e.mobileService.GetProfile(c.Context(), mobileService.GetProfileIn{
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

package tenant_management

import (
	"fmt"
	centralService "hroost/domain/central/tenant/service"
	"hroost/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	CentralService *centralService.Service
}

type TenantManagement struct {
	centralService *centralService.Service
}

func NewTenantManagement(cfg *Config) (*TenantManagement, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for tenant_management required")
	}

	return &TenantManagement{}, nil
}

func (t TenantManagement) CreateTenant(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.CentralAppID:
		var body centralService.CreateTenantIn
		if err := c.BodyParser(&body); err != nil {
			res.Message = err.Error()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(res)
		}

		// call the service
		serviceOut := t.centralService.CreateTenant(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
			res.Data = serviceOut
		} else if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Error = serviceOut.GetError()
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

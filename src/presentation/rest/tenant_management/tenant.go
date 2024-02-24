package tenant_management

import (
	"fmt"

	centralDb "hroost/central/domain/tenant/db"
	centralQueue "hroost/central/domain/tenant/queue"
	centralService "hroost/central/domain/tenant/service"

	"hroost/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

type Central struct {
	Db    *centralDb.Db
	Queue *centralQueue.Queue
}

type Config struct {
	Central *Central
}

type TenantManagement struct {
	central *Central
}

func NewTenantManagement(cfg *Config) (*TenantManagement, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for tenant_management required")
	}
	if cfg.Central == nil {
		return nil, fmt.Errorf("Central required")
	}

	return &TenantManagement{
		central: cfg.Central,
	}, nil
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

		service := centralService.CreateTenant{
			Db:    t.central.Db,
			Queue: t.central.Queue,
		}

		// call the service
		serviceOut := service.Exec(c.Context(), body)

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

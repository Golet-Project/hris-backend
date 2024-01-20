package region

import (
	"fmt"
	"hroost/shared/primitive"

	"hroost/shared/domain/region/service"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Service *service.Service
}

type Region struct {
	service *service.Service
}

func NewRegion(cfg *Config) (*Region, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for region required")
	}

	return &Region{
		service: cfg.Service,
	}, nil
}

func (r Region) FindAllProvince(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	// appId := c.Locals("AppID").(primitive.AppID)

	serviceOut := r.service.FindAllProvince(c.Context())

	res.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
		res.Data = nil
	} else {
		for _, province := range serviceOut.Provinces {
			res.Data = append(res.Data, province)
		}
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(res)
}

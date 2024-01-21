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

func (r Region) FindAllRegencyByProvinceId(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	var req service.FindAllRegencyByProvinceIdIn
	if err := c.QueryParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest)
		res.Message = err.Error()
		return c.JSON(res)
	}

	serviceOut := r.service.FindAllRegencyByProvinceId(c.Context(), req)

	res.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		data := make([]interface{}, len(serviceOut.Regency))
		for i, v := range serviceOut.Regency {
			data[i] = v
		}

		res.Data = data
	} else {
		res.Data = nil
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(res)
}

func (r Region) FindAllDistrictByRegencyId(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	var req service.FindAllDistrictByRegencyIdIn
	if err := c.QueryParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest)
		res.Message = err.Error()
		return c.JSON(res)
	}

	serviceOut := r.service.FindAllDistrictByRegencyId(c.Context(), req)

	res.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		data := make([]interface{}, len(serviceOut.District))
		for i, v := range serviceOut.District {
			data[i] = v
		}

		res.Data = data
	} else {
		res.Data = nil
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(res)
}

func (r Region) FindAllVillageByDistrictId(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	var req service.FindAllVillageByDistrictIdIn
	if err := c.QueryParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest)
		res.Message = err.Error()
		return c.JSON(res)
	}

	serviceOut := r.service.FindAllVillageByDistrictId(c.Context(), req)

	res.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		data := make([]interface{}, len(serviceOut.Village))
		for i, v := range serviceOut.Village {
			data[i] = v
		}

		res.Data = data
	} else {
		res.Data = nil
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(res)
}

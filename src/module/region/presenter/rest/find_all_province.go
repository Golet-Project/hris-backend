package rest

import (
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (r *RegionPresenter) FindAllProvince(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	// appId := c.Locals("AppID").(primitive.AppID)

	serviceOut := r.RegionService.FindAllProvince(c.Context())

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

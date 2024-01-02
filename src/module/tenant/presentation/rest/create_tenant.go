package rest

import (
	"hris/module/shared/primitive"
	"hris/module/tenant/central"

	"github.com/gofiber/fiber/v2"
)

func (t *TenantPresentation) CreateTenant(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.CentralAppID:
		var body central.CreateTenantIn
		if err := c.BodyParser(&body); err != nil {
			res.Message = err.Error()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(res)
		}

		// call the service
		serviceOut := t.central.CreateTenant(c.Context(), body)

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

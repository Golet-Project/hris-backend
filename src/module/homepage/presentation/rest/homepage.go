package rest

import (
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (r *Rest) HomePage(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.MobileAppID:
		c.Status(fiber.StatusOK)
		return c.JSON("OK")
	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

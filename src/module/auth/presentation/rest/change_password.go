package rest

import (
	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func (a *AuthPresentation) InternalChangePassword(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := utils.CopyString(c.Get("X-Cid"))
	token := utils.CopyString(c.Get("X-Api-Key"))

	switch primitive.AppID(appId) {
	case primitive.TenantAppID:
		fallthrough
	case primitive.InternalAppID:
		var body internal.ChangePasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		var serviceOut = a.internal.ChangePassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = serviceOut.GetError()
		}

		c.Status(serviceOut.GetCode())

		return c.JSON(res)

	case primitive.MobileAppID:
		var body mobile.ChangePasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		var serviceOut = a.mobile.ChangePassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = serviceOut.GetError()
		}
		c.Status(serviceOut.GetCode())

		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

}

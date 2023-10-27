package rest

import (
	"hris/module/auth/internal"
	"hris/module/auth/mobile"
	"hris/module/auth/tenant"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

// Handle login request for mobile and web
func (p AuthPresentation) BasicAuthLogin(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		domain := c.Locals("domain").(string)
		var body tenant.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}
		body.Domain = domain

		var loginOut = p.Tenant.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)
	case primitive.InternalAppID:
		var body internal.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = p.Internal.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)

	case primitive.MobileAppID:
		var body mobile.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = p.Mobile.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

}

package rest

import (
	"hris/module/auth/service"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

// Handle login request for mobile and web
func (p AuthPresenter) BasicAuthLogin(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.WebAppID:
		fallthrough
	case primitive.InternalAppID:
		var body service.InternalBasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = p.InternalAuthService.BasicAuthLogin(c.Context(), body)

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

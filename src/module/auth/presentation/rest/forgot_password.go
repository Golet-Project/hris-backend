package rest

import (
	"hris/module/auth/service"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (a *AuthPresenter) ForgotPassword(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.InternalAppID:
		var body service.InternalForgotPasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.AppID = appId

		serviceOut := a.InternalAuthService.ForgotPassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

}

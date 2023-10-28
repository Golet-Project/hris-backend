package rest

import (
	"hris/module/auth/internal"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (a *AuthPresentation) InternalPasswordRecoveryCheck(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	token := c.Get("X-Api-Key")

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.InternalAppID:
		// call the service
		var body internal.PasswordRecoveryTokenCheckIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		serviceOut := a.internal.PasswordRecoveryTokenCheck(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = nil
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

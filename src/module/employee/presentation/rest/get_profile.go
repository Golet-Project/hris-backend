package rest

import (
	"hris/module/employee/mobile"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (e *EmployeePresentation) GetProfile(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.MobileAppID:
		claims := c.Locals("user_auth").(jwt.CustomClaims)

		// call the service
		serviceOut := e.mobile.GetProfile(c.Context(), mobile.GetProfileIn{
			UID: claims.UserUID,
		})

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() <= 500 {
			res.Data = nil
		} else {
			res.Data = serviceOut
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

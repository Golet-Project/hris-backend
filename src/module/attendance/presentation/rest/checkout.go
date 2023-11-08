package rest

import (
	"hris/module/attendance/mobile"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (a AttandancePresentation) Checkout(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	claims := c.Locals("user_auth").(jwt.CustomClaims)

	switch appId {
	case primitive.MobileAppID:
		req := mobile.CheckoutIn{
			UID: claims.UserUID,
		}

		serviceOut := a.mobile.Checkout(c.Context(), req)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Error = serviceOut.GetError()
			res.Data = nil
		} else {
			res.Data = serviceOut.GetMessage()
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

}

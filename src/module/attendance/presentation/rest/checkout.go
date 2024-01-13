package rest

import (
	"hroost/module/attendance/mobile"
	"hroost/module/shared/jwt"
	"hroost/module/shared/primitive"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func (a AttandancePresentation) Checkout(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	claims := c.Locals("user_auth").(jwt.CustomClaims)

	switch appId {
	case primitive.MobileAppID:
		tzString := utils.CopyString(c.Get("local_tz"))
		tz, err := strconv.Atoi(tzString)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = "local_tz header is invalid"
			return c.JSON(res)
		}

		req := mobile.CheckoutIn{
			UID:      claims.UserUID,
			Timezone: primitive.Timezone(tz),
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

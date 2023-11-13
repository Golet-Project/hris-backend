package rest

import (
	"hris/module/homepage/mobile"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func (r *Rest) HomePage(c *fiber.Ctx) error {
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

		var req mobile.HomePageIn
		req.UID = claims.UserUID
		req.Timezone = primitive.Timezone(tz)

		serviceOut := r.mobile.HomePage(c.Context(), req)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 {
			res.Data = serviceOut.GetError()
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

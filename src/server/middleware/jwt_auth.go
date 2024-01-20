package middleware

import (
	centralJwt "hroost/central/lib/jwt"
	mobileJwt "hroost/mobile/lib/jwt"
	tenantJwt "hroost/tenant/lib/jwt"

	"hroost/shared/primitive"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ReqHeader struct {
	Authorization string `reqHeader:"Authorization"`
}

func Jwt() fiber.Handler {
	return func(c *fiber.Ctx) error {
		appId := c.Locals("AppID").(primitive.AppID)
		var headers ReqHeader
		// get the header
		err := c.ReqHeaderParser(&headers)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(map[string]interface{}{
				"message": "invalid header",
			})
		}

		if len(headers.Authorization) == 0 {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(map[string]interface{}{
				"message": "authorization is required",
			})
		}

		splitted := strings.Split(headers.Authorization, " ")
		token := splitted[len(splitted)-1]

		switch appId {
		case primitive.TenantAppID:
			// verify the token
			claims, err := tenantJwt.DecodeAccessToken(token)
			if err != nil {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(map[string]interface{}{
					"message": err.Error(),
				})
			}

			// pass the data
			c.Locals("user_auth", claims)

		case primitive.CentralAppID:
			// verify the token
			claims, err := centralJwt.DecodeAccessToken(token)
			if err != nil {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(map[string]interface{}{
					"message": err.Error(),
				})
			}

			// pass the data
			c.Locals("user_auth", claims)

		case primitive.MobileAppID:
			// verify the token
			claims, err := mobileJwt.DecodeAccessToken(token)
			if err != nil {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(map[string]interface{}{
					"message": err.Error(),
				})
			}

			// pass the data
			c.Locals("user_auth", claims)

		}
		return c.Next()
	}
}

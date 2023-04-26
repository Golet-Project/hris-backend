package middleware

import (
	"hris/module/shared/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ReqHeader struct {
	Authorization string `reqHeader:"Authorization"`
}

func Jwt() fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		// verify the token
		claims, err := jwt.DecodeAccessToken(token)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(map[string]interface{}{
				"message": err.Error(),
			})
		}

		// pass the data
		c.Locals("user_auth", claims)

		return c.Next()
	}
}

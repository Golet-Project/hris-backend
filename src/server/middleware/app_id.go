package middleware

import (
	"hroost/shared/primitive"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/mileusna/useragent"
)

// middleware
func AppId(c *fiber.Ctx) error {
	// except for the change password endpoint
	originalUrl := utils.CopyString(c.OriginalURL())
	if strings.HasPrefix(originalUrl, "/auth/password") && c.Method() == "PUT" {
		return c.Next()
	}

	userAgent := utils.CopyString(string(c.Context().UserAgent()))

	ua := useragent.Parse(userAgent)

	switch ua.Name { // from browser
	case useragent.Opera:
		fallthrough
	case useragent.OperaMini:
		fallthrough
	case useragent.OperaTouch:
		fallthrough
	case useragent.Chrome:
		fallthrough
	case useragent.HeadlessChrome:
		fallthrough
	case useragent.Firefox:
		fallthrough
	case useragent.InternetExplorer:
		fallthrough
	case useragent.Safari:
		fallthrough
	case useragent.Edge:
		fallthrough
	case useragent.Vivaldi:
		appId := utils.CopyString(c.Get("X-App-ID"))
		domain := utils.CopyString(c.Get("X-Domain"))
		switch appId {
		case primitive.CentralAppID.String():
			c.Locals("AppID", primitive.CentralAppID)
			return c.Next()
		case primitive.TenantAppID.String():
			c.Locals("AppID", primitive.TenantAppID)
			c.Locals("domain", domain)
			return c.Next()

		default:
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(map[string]string{
				"message": "invalid app ID",
			})
		}

	default: // verify from mobile devices
		switch ua.OS {
		case useragent.Android:
			fallthrough
		case useragent.IOS:
			c.Locals("AppID", primitive.MobileAppID)
			return c.Next()

		default:
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(map[string]string{
				"message": "invalid client",
			})
		}
	}
}

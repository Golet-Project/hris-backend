package homepage

import (
	"fmt"
	"hroost/shared/primitive"

	mobileDb "hroost/mobile/domain/homepage/db"
	mobileService "hroost/mobile/domain/homepage/service"
	mobileJwt "hroost/mobile/lib/jwt"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Mobile struct {
	Db *mobileDb.Db
}

type Config struct {
	Mobile *Mobile
}

type HomePage struct {
	mobile *Mobile
}

func NewHomepage(cfg *Config) (*HomePage, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for homepage required")
	}
	if cfg.Mobile == nil {
		return nil, fmt.Errorf("mobile module required")
	}

	return &HomePage{
		mobile: cfg.Mobile,
	}, nil
}

func (h HomePage) HomePage(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.MobileAppID:
		claims := c.Locals("user_auth").(mobileJwt.CustomClaims)
		tzString := utils.CopyString(c.Get("local_tz"))
		tz, err := strconv.Atoi(tzString)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = "local_tz header is invalid"
			return c.JSON(res)
		}

		var req mobileService.HomePageIn
		req.UID = claims.UserUID
		req.Timezone = primitive.Timezone(tz)

		service := mobileService.HomePage{
			Db: h.mobile.Db,
		}

		serviceOut := service.Exec(c.Context(), req)

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

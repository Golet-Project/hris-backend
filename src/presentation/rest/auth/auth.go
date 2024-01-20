package auth

import (
	"fmt"
	centralService "hroost/central/domain/auth/service"
	mobileService "hroost/mobile/domain/auth/service"
	"hroost/shared/primitive"
	tenantService "hroost/tenant/domain/auth/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Config struct {
	CentralService *centralService.Service
	MobileService  *mobileService.Service
	TenantService  *tenantService.Service
}

type Auth struct {
	centralService *centralService.Service
	mobileService  *mobileService.Service
	tenantService  *tenantService.Service
}

func NewAuth(cfg *Config) (*Auth, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for auth required")
	}

	return &Auth{
		centralService: cfg.CentralService,
		mobileService:  cfg.MobileService,
		tenantService:  cfg.TenantService,
	}, nil
}

func (a Auth) BasicAuthLogin(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		var body tenantService.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = a.tenantService.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)
	case primitive.CentralAppID:
		var body centralService.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = a.centralService.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)

	case primitive.MobileAppID:
		var body mobileService.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		var loginOut = a.mobileService.BasicAuthLogin(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		c.Status(loginOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

func (a Auth) ChangePassword(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := utils.CopyString(c.Get("X-Cid"))
	token := utils.CopyString(c.Get("X-Api-Key"))

	switch primitive.AppID(appId) {
	case primitive.TenantAppID:
		fallthrough
	case primitive.CentralAppID:
		var body centralService.ChangePasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		var serviceOut = a.centralService.ChangePassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = serviceOut.GetError()
		}

		c.Status(serviceOut.GetCode())

		return c.JSON(res)

	case primitive.MobileAppID:
		var body mobileService.ChangePasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		var serviceOut = a.mobileService.ChangePassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = serviceOut.GetError()
		}
		c.Status(serviceOut.GetCode())

		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

}

func (a Auth) ForgotPassword(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.CentralAppID:
		var body centralService.ForgotPasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.AppID = appId

		serviceOut := a.centralService.ForgotPassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
			res.Data = serviceOut.GetMessage()
		} else if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Error = serviceOut.GetError()
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	case primitive.MobileAppID:
		var body mobileService.ForgotPasswordIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		serviceOut := a.mobileService.ForgotPassword(c.Context(), body)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
			res.Data = serviceOut.GetMessage()
		} else if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Error = serviceOut.GetError()
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

func (a Auth) OAuthLogin(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.CentralAppID:
		// call the service
		serviceOut := a.centralService.OAuthLogin(c.Context())
		response.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
			response.Data = serviceOut
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(response)

	default:
		response.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response)
	}
}

func (a Auth) OAuthCallback(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	appId := c.Get("X-App-ID")

	switch appId {
	case primitive.CentralAppID.String():
		var query centralService.OAuthCallbackIn
		c.QueryParser(&query)

		// call the service
		serviceOut := a.centralService.OAuthCallback(c.Context(), query)
		response.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
			response.Data = serviceOut
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(response)

	default:
		response.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response)
	}
}

func (a Auth) PasswordRecoveryCheck(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	token := c.Get("X-Api-Key")

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.CentralAppID:
		// call the service
		var body centralService.PasswordRecoveryTokenCheckIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.Token = token

		serviceOut := a.centralService.PasswordRecoveryTokenCheck(c.Context(), body)

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

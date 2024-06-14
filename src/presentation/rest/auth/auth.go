package auth

import (
	"fmt"
	"hroost/shared/primitive"

	tenantDb "hroost/tenant/domain/auth/db"
	tenantService "hroost/tenant/domain/auth/service"

	centralDb "hroost/central/domain/auth/db"
	centralMemory "hroost/central/domain/auth/memory"
	centralService "hroost/central/domain/auth/service"

	mobileDb "hroost/mobile/domain/auth/db"
	mobileMemory "hroost/mobile/domain/auth/memory"
	mobileService "hroost/mobile/domain/auth/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type Central struct {
	Db           *centralDb.Db
	Memory       *centralMemory.Memory
	Oauth2Google *oauth2.Config
}

type Mobile struct {
	Db     *mobileDb.Db
	Memory *mobileMemory.Memory
}

type Tenant struct {
	Db *tenantDb.Db
}

type Config struct {
	Central *Central
	Mobile  *Mobile
	Tenant  *Tenant
}

type Auth struct {
	central *Central
	mobile  *Mobile
	tenant  *Tenant
}

func NewAuth(cfg *Config) (*Auth, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dependency for auth required")
	}

	return &Auth{
		central: cfg.Central,
		mobile:  cfg.Mobile,
		tenant:  cfg.Tenant,
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

		service := tenantService.BasicAuthLogin{
			Db: a.tenant.Db,
		}

		loginOut := service.Exec(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut
		} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
			res.Error = loginOut.GetError()
		}

		// NOTE: currently we only set access token in the header for tenant
		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = loginOut.AccessToken
		cookie.MaxAge = 24 * 3600
		cookie.HTTPOnly = true

		c.Cookie(cookie)
		c.Status(loginOut.GetCode())
		return c.JSON(res)
	case primitive.CentralAppID:
		var body centralService.BasicAuthLoginIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		service := centralService.BasicAuthLogin{
			Db: a.central.Db,
		}

		loginOut := service.Exec(c.Context(), body)

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

		service := mobileService.BasicAuthLogin{
			Db: a.mobile.Db,
		}

		loginOut := service.Exec(c.Context(), body)

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

		service := centralService.ChangePassword{
			Memory:               a.central.Memory,
			Db:                   a.central.Db,
			GenerateFromPassword: bcrypt.GenerateFromPassword,
		}

		serviceOut := service.Exec(c.Context(), body)

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

		service := mobileService.ChangePassword{
			Db:     a.mobile.Db,
			Memory: a.mobile.Memory,
		}

		serviceOut := service.Exec(c.Context(), body)

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

		service := centralService.ForgotPassword{
			Db:     a.central.Db,
			Memory: a.central.Memory,
		}

		serviceOut := service.Exec(c.Context(), body)

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

		service := mobileService.ForgotPassword{
			Db:     a.mobile.Db,
			Memory: a.mobile.Memory,
		}

		serviceOut := service.Exec(c.Context(), body)

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

		service := centralService.OAuthLogin{
			OAuth2Google: a.central.Oauth2Google,
		}
		// call the service
		serviceOut := service.Exec(c.Context())
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

		service := centralService.OAuthCallback{
			Db:           a.central.Db,
			Oauth2Google: a.central.Oauth2Google,
		}

		// call the service
		serviceOut := service.Exec(c.Context(), query)
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

		service := centralService.PasswordRecoveryTokenCheck{
			Memory: a.central.Memory,
		}

		serviceOut := service.Exec(c.Context(), body)

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

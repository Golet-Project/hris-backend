package rest

import (
	"hris/module/auth/central"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

// Handle incoming request to perform o-auth login
func (p AuthPresentation) OAuthLogin(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		fallthrough
	case primitive.CentralAppID:
		// call the service
		serviceOut := p.central.OAuthLogin(c.Context())
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

func (p AuthPresentation) OAuthCallback(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	appId := c.Get("X-App-ID")

	switch appId {
	case primitive.CentralAppID.String():
		var query central.OAuthCallbackIn
		c.QueryParser(&query)

		// call the service
		serviceOut := p.central.OAuthCallback(c.Context(), query)
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

package rest

import (
	"hris/module/auth/service"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

// Handle incoming request to perform o-auth login
func (p AuthPresenter) OAuthLogin(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	// call the service
	serviceOut := p.AuthService.OAuthLogin(c.Context())
	response.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		response.Data = serviceOut
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(response)
}

func (p AuthPresenter) OAuthCallback(c *fiber.Ctx) error {
	var response primitive.BaseResponse

	var query service.OAuthCallbackIn
	c.QueryParser(&query)

	// call the service
	serviceOut := p.AuthService.OAuthCallback(c.Context(), query)
	response.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		response.Data = serviceOut
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(response)
}

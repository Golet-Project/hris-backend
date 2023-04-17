package rest

import (
	"hris/module/auth/service"
	"hris/module/shared/primitive"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Handle login request for mobile and web
func (p AuthPresenter) BasicAuthLogin(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	var body service.LoginIn
	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		res.Message = err.Error()
		return c.JSON(res)
	}

	var loginOut = p.AuthService.BasicAuthLogin(c.Context(), body)

	res.Message = loginOut.GetMessage()

	if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
		res.Data = loginOut
	} else if loginOut.GetCode() >= 400 && loginOut.GetCode() < 500 {
		res.Data = loginOut.GetError()
	}

	c.Status(loginOut.GetCode())
	return c.JSON(res)
}

// Handle incomin request for auth callback when using one-tap-sign
func (p AuthPresenter) OneTapCallback(c *fiber.Ctx) error {
	var response primitive.BaseResponse
	var in service.AuthCallbackIn

	// bind the cookie
	in.CookieToken = c.Cookies("g_csrf_token")
	// bind the body
	if err := c.BodyParser(&in); err != nil {
		c.Status(http.StatusBadRequest)
		response.Message = err.Error()
		return c.JSON(response)
	}

	// call the service
	out := p.AuthService.OneTapCallback(c.Context(), in)
	response.Message = out.GetMessage()

	if out.GetCode() >= 200 && out.GetCode() < 400 {
		response.Data = out
	} else if out.GetCode() >= 400 && out.GetCode() < 500 {
		response.Data = out.GetError()
	}

	c.Status(out.GetCode())
	return c.JSON(response)
}

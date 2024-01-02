package rest

import (
	"hris/module/employee/tenant"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (e *EmployeePresentation) CreateEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(jwt.TenantCustomClaims)

		var req tenant.FindAllEmployeeIn
		req.Domain = claims.Domain

		// call the services
		serviceOut := e.tenant.FindAllEmployee(c.Context(), req)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 {
			res.Data = nil
		} else {
			data := make([]interface{}, len(serviceOut.Employee))
			for i, v := range serviceOut.Employee {
				data[i] = v
			}
			res.Data = data
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

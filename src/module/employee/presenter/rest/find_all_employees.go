package rest

import (
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (e *EmployeePresentation) FindAllEmployee(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	appId := c.Locals("AppID").(primitive.AppID)
	switch appId {
	case primitive.TenantAppID:
		// call the service
		serviceOut := e.Tenant.FindAllEmployee(c.Context())

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Data = nil
		} else {
			for _, employee := range serviceOut.Employee {
				res.Data = append(res.Data, employee)
			}
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

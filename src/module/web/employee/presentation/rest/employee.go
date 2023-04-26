package rest

import (
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (p EmployeePresenter) FindAllEmployee(c *fiber.Ctx) error {
	var response primitive.BaseResponseArray

	// call the service
	serviceOut := p.EmployeeService.FindAllEmployee(c.Context())
	response.Message = serviceOut.GetMessage()

	if serviceOut.GetCode() >= 200 && serviceOut.GetCode() < 400 {
		for _, val := range serviceOut.Data {
			response.Data = append(response.Data, val)
		}
	}

	c.Status(serviceOut.GetCode())
	return c.JSON(response)
}

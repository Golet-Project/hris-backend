package rest

import (
	"hris/module/attendance/tenant"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"

	"github.com/gofiber/fiber/v2"
)

func (a AttandancePresentation) FindAllAttendance(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(jwt.TenantCustomClaims)

		serviceOut := a.tenant.FindAllAttendance(c.Context(), tenant.FindAllAttendanceIn{
			Domain: claims.Domain,
		})

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 {
			res.Data = nil
		} else {
			data := make([]interface{}, len(serviceOut.Attendance))
			for i, v := range serviceOut.Attendance {
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

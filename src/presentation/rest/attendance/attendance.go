package attendance

import (
	"fmt"
	mobileJwt "hroost/mobile/lib/jwt"
	tenantJwt "hroost/tenant/lib/jwt"

	"hroost/shared/primitive"
	"strconv"

	mobileDb "hroost/mobile/domain/attendance/db"
	mobileService "hroost/mobile/domain/attendance/service"

	tenantDb "hroost/tenant/domain/attendance/db"
	tenantService "hroost/tenant/domain/attendance/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Mobile struct {
	Db *mobileDb.Db
}

type Tenant struct {
	Db *tenantDb.Db
}

type Config struct {
	Mobile *Mobile
	Tenant *Tenant
}

type Attendance struct {
	mobile *Mobile
	tenant *Tenant
}

func NewAttendance(cfg *Config) (*Attendance, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config for attendance required")
	}

	return &Attendance{
		mobile: cfg.Mobile,
		tenant: cfg.Tenant,
	}, nil
}

func (a Attendance) AddAttendance(c *fiber.Ctx) error {
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

		var body mobileService.AddAttendanceIn
		if err := c.BodyParser(&body); err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = err.Error()
			return c.JSON(res)
		}

		body.UID = claims.UserUID
		body.Timezone = primitive.Timezone(tz)

		service := mobileService.AddAttendance{
			Db: a.mobile.Db,
		}

		var loginOut = service.Exec(c.Context(), body)

		res.Message = loginOut.GetMessage()

		if loginOut.GetCode() >= 200 && loginOut.GetCode() < 400 {
			res.Data = loginOut.GetMessage()
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

func (a Attendance) Checkout(c *fiber.Ctx) error {
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

		req := mobileService.CheckoutIn{
			UID:      claims.UserUID,
			Timezone: primitive.Timezone(tz),
		}

		service := mobileService.Checkout{
			Db: a.mobile.Db,
		}

		serviceOut := service.Exec(c.Context(), req)

		res.Message = serviceOut.GetMessage()

		if serviceOut.GetCode() >= 400 && serviceOut.GetCode() < 500 {
			res.Error = serviceOut.GetError()
			res.Data = nil
		} else {
			res.Data = serviceOut.GetMessage()
		}

		c.Status(serviceOut.GetCode())
		return c.JSON(res)

	default:
		res.Message = "invalid app id"
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
}

func (a Attendance) FindAllAttendance(c *fiber.Ctx) error {
	var res primitive.BaseResponseArray

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.TenantAppID:
		claims := c.Locals("user_auth").(tenantJwt.CustomClaims)

		service := tenantService.FindAllAttendance{
			Db: a.tenant.Db,
		}

		serviceOut := service.Exec(c.Context(), tenantService.FindAllAttendanceIn{
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

func (a Attendance) GetTodayAttendance(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)

	switch appId {
	case primitive.MobileAppID:
		tzString := utils.CopyString(c.Get("local_tz"))
		tz, err := strconv.Atoi(tzString)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			res.Message = "local_tz header is invalid"
			return c.JSON(res)
		}

		claims := c.Locals("user_auth").(mobileJwt.CustomClaims)

		var req mobileService.GetTodayAttendanceIn
		req.EmployeeUID = claims.UserUID
		req.Timezone = primitive.Timezone(tz)

		service := mobileService.GetTodayAttendance{
			Db: a.mobile.Db,
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

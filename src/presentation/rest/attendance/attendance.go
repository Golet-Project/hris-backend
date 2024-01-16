package attendance

import (
	"fmt"
	"hroost/module/shared/jwt"
	"hroost/module/shared/primitive"
	"strconv"

	mobileService "hroost/domain/mobile/attendance/service"
	tenantService "hroost/domain/tenant/attendance/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Config struct {
	MobileService *mobileService.Service
	TenantService *tenantService.Service
}

type Attendance struct {
	mobileService *mobileService.Service
	tenantService *tenantService.Service
}

func NewAttendance(cfg *Config) (*Attendance, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config for attendance required")
	}

	return &Attendance{
		mobileService: cfg.MobileService,
		tenantService: cfg.TenantService,
	}, nil
}

func (a Attendance) AddAttendance(c *fiber.Ctx) error {
	var res primitive.BaseResponse

	appId := c.Locals("AppID").(primitive.AppID)
	claims := c.Locals("user_auth").(jwt.CustomClaims)

	switch appId {
	case primitive.MobileAppID:
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
		var loginOut = a.mobileService.AddAttendance(c.Context(), body)

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
	claims := c.Locals("user_auth").(jwt.CustomClaims)

	switch appId {
	case primitive.MobileAppID:
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

		serviceOut := a.mobileService.Checkout(c.Context(), req)

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
		claims := c.Locals("user_auth").(jwt.TenantCustomClaims)

		serviceOut := a.tenantService.FindAllAttendance(c.Context(), tenantService.FindAllAttendanceIn{
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

		claims := c.Locals("user_auth").(jwt.CustomClaims)

		var req mobileService.GetTodayAttendanceIn
		req.EmployeeUID = claims.UserUID
		req.Timezone = primitive.Timezone(tz)

		serviceOut := a.mobileService.GetTodayAttendance(c.Context(), req)

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

package attendance

import (
	"hris/module/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func (a Attendance) Route(app *fiber.App) {
	app.Post("/attendance", middleware.Jwt(), a.AttendancePresentation.AddAttendance)
	app.Put("/attendance", middleware.Jwt(), a.AttendancePresentation.Checkout)
}

package employee

import (
	"hris/module/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func (e Employee) Route(app *fiber.App) {
	app.Get("/employee", e.EmployeePresentation.FindAllEmployee)
	app.Get("/profile", middleware.Jwt(), e.EmployeePresentation.GetProfile)
}

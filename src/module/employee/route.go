package employee

import (
	"hroost/module/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func (e Employee) Route(app *fiber.App) {
	app.Get("/employee", middleware.Jwt(), e.EmployeePresentation.FindAllEmployee)
	app.Post("/employee", middleware.Jwt(), e.EmployeePresentation.CreateEmployee)

	app.Get("/profile", middleware.Jwt(), e.EmployeePresentation.GetProfile)
}

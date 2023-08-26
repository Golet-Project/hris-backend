package employee

import "github.com/gofiber/fiber/v2"

func (e Employee) Route(app *fiber.App) {
	app.Get("/employees", e.EmployeePresenter.FindAllEmployees)
}

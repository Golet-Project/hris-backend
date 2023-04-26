package employee

import (
	"hris/module/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func (e Employee) Route(r fiber.Router) {
	r.Use(middleware.Jwt())

	r.Get("/employees", e.EmployeePresenter.FindAllEmployee)
}

package homepage

import (
	"hris/module/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

func (h HomePage) Route(app *fiber.App) {
	app.Get("/homepage", middleware.Jwt(), h.Rest.HomePage)
}

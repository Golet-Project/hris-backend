package tenant

import "github.com/gofiber/fiber/v2"

func (t Tenant) Route(app *fiber.App) {
	app.Post("/tenant", t.Internal_TenantPresenter.CreateTenant)
}

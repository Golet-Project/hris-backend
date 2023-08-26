package region

import "github.com/gofiber/fiber/v2"

func (r Region) Route(app *fiber.App) {
	app.Get("/provinces", r.RegionPresenter.FindAllProvince)
}

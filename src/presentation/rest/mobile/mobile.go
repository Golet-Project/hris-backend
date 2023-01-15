package mobile

import (
	"hris/presentation/rest/mobile/user"

	"github.com/gofiber/fiber/v2"
)

type Dependency struct {
	UserPresenter *user.UserPresenter
}

func (d Dependency) Route(r fiber.Router) {
	userGroup := r.Group("/user")
	userGroup.Get("/", d.UserPresenter.SayHello)
}

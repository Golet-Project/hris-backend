package user

import (
	"hris/module/mobile/user"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserPresenter struct {
	UserService user.Service
}

func (u *UserPresenter) SayHello(c *fiber.Ctx) error {
	c.Status(http.StatusOK)
	return c.JSON(u.UserService.SayHello())
}

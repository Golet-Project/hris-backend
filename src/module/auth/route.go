package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (a Auth) Route(app *fiber.App) {
	app.Post("/auth/login", a.AuthPresenter.BasicAuthLogin)

	app.Post("/auth/callback", a.AuthPresenter.OneTapCallback)

	app.Post("/oauth/google/login", a.AuthPresenter.OAuthLogin)
	app.Get("/oauth/google/callback", a.AuthPresenter.OAuthCallback)
}

package rest

import (
	"hris/module/auth/presentation/rest"

	"github.com/gofiber/fiber/v2"
)

type Dependency struct {
	AuthPresenter *rest.AuthPresenter
}

func (d Dependency) Route(app *fiber.App) {
	app.Post("/auth/login", d.AuthPresenter.Login)
	app.Post("/auth/callback", d.AuthPresenter.AuthCallback)

	app.Post("/oauth/google/login", d.AuthPresenter.OAuthLogin)
	app.Get("/oauth/google/callback", d.AuthPresenter.OAuthCallback)
}

package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (a Auth) Route(app *fiber.App) {
	app.Post("/auth/login", a.AuthPresenter.BasicAuthLogin)

	app.Post("/oauth/google/login", a.AuthPresenter.OAuthLogin)
	app.Get("/oauth/google/callback", a.AuthPresenter.OAuthCallback)

	app.Post("/auth/forgot-password", a.AuthPresenter.ForgotPassword)
	app.Post("/auth/password-recovery/check", a.AuthPresenter.InternalPasswordRecoveryCheck)
	app.Put("/auth/password", a.AuthPresenter.InternalChangePassword)

	// user
	// app.Post("/auth/login")
}

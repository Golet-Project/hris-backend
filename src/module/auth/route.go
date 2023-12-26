package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (a Auth) Route(app *fiber.App) {
	app.Post("/auth/login", a.AuthPresentation.BasicAuthLogin)

	app.Post("/oauth/google/login", a.AuthPresentation.OAuthLogin)
	app.Get("/oauth/google/callback", a.AuthPresentation.OAuthCallback)

	app.Post("/auth/forgot-password", a.AuthPresentation.ForgotPassword)
	app.Post("/auth/password-recovery/check", a.AuthPresentation.CentralPasswordRecoveryCheck)
	app.Put("/auth/password", a.AuthPresentation.CentralChangePassword)

	// user
	// app.Post("/auth/login")
}

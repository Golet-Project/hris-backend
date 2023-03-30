package service

import (
	"context"
	"errors"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"
	"net/http"

	"github.com/jackc/pgx/v5"
	"google.golang.org/api/idtoken"
)

type AuthCallbackIn struct {
	Credential string `json:"credential"`

	CookieToken string `json:"-"`
	GCSRFToken  string `form:"g_csrf_token"`
}

type AuthCallbackOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

// Handle request callback if login using google one-tap-sign
func (s *AuthService) AuthCallback(ctx context.Context, in AuthCallbackIn) (out AuthCallbackOut) {
	if !csrfTokenValid(in.CookieToken, in.GCSRFToken) {
		out.SetResponse(http.StatusBadRequest, "Invalid CSRF token")
		return
	}

	// validate the id token
	payload, err := idtoken.Validate(ctx, in.Credential, s.oauth2Cfg.ClientID)
	if err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid ID token")
		return
	}

	if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
		out.SetResponse(http.StatusUnauthorized, "invalid ID token")
		return
	}

	userEmail, ok := payload.Claims["email"].(string)
	if !ok {
		out.SetResponse(http.StatusUnauthorized, "invalid ID token")
		return
	}

	// get the user data
	user, err := s.AuthRepo.GetLoginCredential(ctx, userEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusUnauthorized, "user not registered")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// generate jwt token
	out.AccessToken = jwt.GenerateAccessToken(user.Email)
	out.SetResponse(http.StatusOK, "success")
	return

}

func csrfTokenValid(cookieToken, bodyToken string) bool {
	if cookieToken == "" || bodyToken == "" {
		return false
	}

	if cookieToken != bodyToken {
		return false
	}

	return true
}

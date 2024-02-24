package service

import (
	"context"
	"hroost/shared/primitive"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthLoginOut struct {
	primitive.CommonResult

	Url string `json:"url"`
}

type OAuth2Google interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
}

type OAuthLogin struct {
	OAuth2Google OAuth2Google
}

// OAuthLogin redirect user to Google's consent page to ask for permission
func (s *OAuthLogin) Exec(ctx context.Context) (out OAuthLoginOut) {
	out.Url = s.OAuth2Google.AuthCodeURL(oauthState)
	out.SetResponse(http.StatusTemporaryRedirect, "Redirecting to Google")
	return
}

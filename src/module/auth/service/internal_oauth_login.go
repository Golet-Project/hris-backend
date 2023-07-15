package service

import (
	"context"
	"hris/module/shared/primitive"
	"net/http"
)

type OAuthLoginOut struct {
	primitive.CommonResult

	Url string `json:"url"`
}

func (s *InternalAuthService) OAuthLogin(ctx context.Context) (out OAuthLoginOut) {
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	out.Url = s.oauth2Cfg.AuthCodeURL(oauthState)
	out.SetResponse(http.StatusTemporaryRedirect, "Redirecting to Google")
	return
}

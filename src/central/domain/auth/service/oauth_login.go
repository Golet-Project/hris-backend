package service

import (
	"context"
	"hroost/shared/primitive"
	"net/http"
)

type OAuthLoginOut struct {
	primitive.CommonResult

	Url string `json:"url"`
}

// OAuthLogin redirect user to Google's consent page to ask for permission
func (s *Service) OAuthLogin(ctx context.Context) (out OAuthLoginOut) {
	out.Url = s.oauth2Cfg.AuthCodeURL(oauthState)
	out.SetResponse(http.StatusTemporaryRedirect, "Redirecting to Google")
	return
}

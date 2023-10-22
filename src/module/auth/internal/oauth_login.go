package internal

import (
	"context"
	"hris/module/shared/primitive"
	"net/http"
)

type OAuthLoginOut struct {
	primitive.CommonResult

	Url string `json:"url"`
}

// OAuthLogin redirect user to Google's consent page to ask for permission
func (i *Internal) OAuthLogin(ctx context.Context) (out OAuthLoginOut) {
	out.Url = i.oauth2Cfg.AuthCodeURL(oauthState)
	out.SetResponse(http.StatusTemporaryRedirect, "Redirecting to Google")
	return
}

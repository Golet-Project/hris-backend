package service

import (
	"hris/module/auth/repo/auth"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type AuthService struct {
	AuthRepo *auth.Repository

	oauth2Cfg *oauth2.Config
}

var oauthState = os.Getenv("OAUTH_STATE")

func NewAuthService(authRepo *auth.Repository) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,

		oauth2Cfg: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			Endpoint:     endpoints.Google,
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
	}
}

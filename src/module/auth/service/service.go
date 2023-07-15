package service

import (
	"hris/module/auth/repo/auth"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type InternalAuthService struct {
	AuthRepo *auth.Repository

	oauth2Cfg *oauth2.Config
}

var oauthState = os.Getenv("OAUTH_STATE")

func NewInternalAuthService(authRepo *auth.Repository) *InternalAuthService {
	return &InternalAuthService{
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

type MobileAuthService struct {
	AuthRepo *auth.Repository
}

func NewMobileAuthService(authRepo *auth.Repository) *MobileAuthService {
	return &MobileAuthService{
		AuthRepo: authRepo,
	}
}

type WebAuthService struct {
	AuthRepo *auth.Repository

	oauth2Cfg *oauth2.Config
}

func NewWebAuthService(authRepo *auth.Repository) *WebAuthService {
	return &WebAuthService{
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
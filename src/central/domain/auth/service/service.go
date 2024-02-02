package service

import (
	"fmt"
	"hroost/central/domain/auth/db"
	"hroost/central/domain/auth/memory"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var oauthState = os.Getenv("OAUTH_STATE")

type Config struct {
	Db     db.IDbStore
	Memory memory.IMemory
}

type Service struct {
	db     db.IDbStore
	memory memory.IMemory

	oauth2Cfg *oauth2.Config
}

func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config for service required")
	}

	return &Service{
		db:     cfg.Db,
		memory: cfg.Memory,

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
	}, nil
}

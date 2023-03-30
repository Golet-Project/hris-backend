package service

import (
	"context"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type OAuthCallbackIn struct {
	Code string `query:"code"`
}

type OAuthCallbackOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

// Handle callback when login usig google OAuth
func (s *AuthService) OAuthCallback(ctx context.Context, query OAuthCallbackIn) (out OAuthCallbackOut) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token, err := s.oauth2Cfg.Exchange(ctx, query.Code)
	log.Println(query.Code)
	if err != nil {
		out.SetResponse(http.StatusUnauthorized, "failed to exchange token")
		return
	}

	// get user info
	oauth2Service, err := googleOAuth2.NewService(ctx,
		option.WithTokenSource(s.oauth2Cfg.TokenSource(ctx, token)),
	)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "failed to create oauth2 service", err)
		return
	}
	userInfo, err := googleOAuth2.NewUserinfoService(oauth2Service).Get().Do()
	if err != nil {
		out.SetResponse(http.StatusUnauthorized, "failed to fetch user info")
		return
	}

	// get the user data
	user, err := s.AuthRepo.GetLoginCredential(ctx, userInfo.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			out.SetResponse(http.StatusUnauthorized, "user not registered")
			return
		}

		out.SetResponse(http.StatusInternalServerError, "failed to get user data")
		return
	}

	// generate jwt token
	out.AccessToken = jwt.GenerateAccessToken(user.Email)

	out.SetResponse(http.StatusOK, "success")
	return
}

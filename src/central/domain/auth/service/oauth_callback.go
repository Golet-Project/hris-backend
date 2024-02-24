package service

import (
	"context"
	"hroost/central/domain/auth/model"
	"hroost/central/lib/jwt"
	"hroost/shared/primitive"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	googleOAuth2 "google.golang.org/api/oauth2/v2"
)

type OAuthCallbackIn struct {
	Code string `json:"code"`
}

type OAuthCallbackOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

type OAuthCallbackDb interface {
	GetLoginCredential(ctx context.Context, email string) (data model.GetLoginCredentialOut, err *primitive.RepoError)
}
type OAuthCallbackGoogle interface {
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource
}

type OAuthCallback struct {
	Db           OAuthCallbackDb
	Oauth2Google OAuthCallbackGoogle
}

// OAuthCallback handle callback when login usig google OAuth
func (s *OAuthCallback) Exec(ctx context.Context, query OAuthCallbackIn) (out OAuthCallbackOut) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token, err := s.Oauth2Google.Exchange(ctx, query.Code)
	if err != nil {
		out.SetResponse(http.StatusUnauthorized, "failed to exchange token")
		return
	}

	// get user info
	oauth2Service, err := googleOAuth2.NewService(ctx,
		option.WithTokenSource(s.Oauth2Google.TokenSource(ctx, token)),
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
	user, repoError := s.Db.GetLoginCredential(ctx, userInfo.Email)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusUnauthorized, "user not registered")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "failed to get user data")
			return
		}
	}

	// generate jwt token
	out.AccessToken = jwt.GenerateAccessToken(user.Email)

	out.SetResponse(http.StatusOK, "success")
	return
}

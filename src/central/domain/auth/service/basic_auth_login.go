package service

import (
	"context"
	"hroost/central/domain/auth/model"
	"hroost/central/lib/jwt"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type BasicAuthLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BasicAuthLoginOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

type BasicAuthLoginDb interface {
	GetLoginCredential(ctx context.Context, email string) (out model.GetLoginCredentialOut, err *primitive.RepoError)
}

type BasicAuthLogin struct {
	Db BasicAuthLoginDb
}

func (s *BasicAuthLogin) Exec(ctx context.Context, body BasicAuthLoginIn) (out BasicAuthLoginOut) {
	// validate the request body
	if err := s.ValidateBasicAuthLoginBody(body); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login data
	loginData, err := s.Db.GetLoginCredential(ctx, body.Email)
	if err != nil {
		switch err.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(loginData.Password.String), []byte(body.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	// generate access token
	out.AccessToken = jwt.GenerateAccessToken(loginData.UserUID)

	out.SetResponse(http.StatusOK, "login success")
	return
}

func (s *BasicAuthLogin) ValidateBasicAuthLoginBody(body BasicAuthLoginIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
	} else if issues := utils.IsEmailValid(body.Email); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	// validate password
	if len(body.Password) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := utils.IsPasswordValid(body.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(body.Password) > 25 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooLong,
			Field:   "password",
			Message: "password must be less than 25 characters",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

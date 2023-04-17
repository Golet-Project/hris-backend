package service

import (
	"context"
	"errors"
	"hris/module/shared/jwt"
	"hris/module/shared/primitive"
	"hris/module/shared/validator"
	"net/http"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func ValidateLoginBody(body LoginIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
	} else if issues := validator.IsEmailValid(body.Email); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	// validate password
	if len(body.Password) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := validator.IsPasswordValid(body.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

type LoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

// Handle login using email and password
func (s *AuthService) BasicAuthLogin(ctx context.Context, reqBody LoginIn) (out LoginOut) {
	// validate the request body
	if err := ValidateLoginBody(reqBody); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login data
	loginData, err := s.AuthRepo.GetLoginCredential(ctx, reqBody.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(loginData.Password), []byte(reqBody.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	// generate access token
	out.AccessToken = jwt.GenerateAccessToken(loginData.UserUID)

	out.SetResponse(http.StatusOK, "login success")
	return
}

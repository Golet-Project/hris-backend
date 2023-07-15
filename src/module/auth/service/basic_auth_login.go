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

type InternalBasicAuthLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InternalBasicAuthLoginOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

func ValidateInternalBasicAuthLoginBody(body InternalBasicAuthLoginIn) *primitive.RequestValidationError {
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

// Handle login using email and password
func (s *InternalAuthService) BasicAuthLogin(ctx context.Context, reqBody InternalBasicAuthLoginIn) (out InternalBasicAuthLoginOut) {
	// validate the request body
	if err := ValidateInternalBasicAuthLoginBody(reqBody); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login data
	loginData, err := s.AuthRepo.InternalGetLoginCredential(ctx, reqBody.Email)
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
	if err := bcrypt.CompareHashAndPassword([]byte(loginData.Password.String), []byte(reqBody.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	// generate access token
	out.AccessToken = jwt.GenerateAccessToken(loginData.UserUID)

	out.SetResponse(http.StatusOK, "login success")
	return
}

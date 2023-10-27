package tenant

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

type BasicAuthLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Domain   string
}

type BasicAuthLoginOut struct {
	primitive.CommonResult

	AccessToken string `json:"access_token"`
}

func ValidateBasicAuthLoginBody(body BasicAuthLoginIn) *primitive.RequestValidationError {
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

	// validate domain
	if len(body.Domain) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "domain",
			Message: "domain header is required",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

func (t *Tenant) BasicAuthLogin(ctx context.Context, body BasicAuthLoginIn) (out BasicAuthLoginOut) {
	// validate request body
	if err := ValidateBasicAuthLoginBody(body); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login credentials
	adminCredential, err := t.db.GetLoginCredential(ctx, body.Email, body.Domain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(adminCredential.Password), []byte(body.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	out.AccessToken = jwt.GenerateAccessToken(adminCredential.UserID)

	return
}
